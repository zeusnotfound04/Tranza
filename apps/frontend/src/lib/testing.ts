// Test utilities for API and component testing
import { apiClient, TranzaAPIClient } from '../services/api';

export interface TestConfig {
  baseURL: string;
  testToken?: string;
  timeout: number;
}

export interface TestResults {
  passed: number;
  failed: number;
  total: number;
  details: TestDetail[];
}

export interface TestDetail {
  name: string;
  status: 'pass' | 'fail' | 'skip';
  message: string;
  duration: number;
  error?: string;
}

// Test runner class
export class IntegrationTester {
  private config: TestConfig;
  private results: TestDetail[] = [];

  constructor(config: TestConfig) {
    this.config = config;
  }

  async runTest(name: string, testFn: () => Promise<void>): Promise<TestDetail> {
    const startTime = Date.now();
    let result: TestDetail;

    try {
      await Promise.race([
        testFn(),
        new Promise((_, reject) => 
          setTimeout(() => reject(new Error('Test timeout')), this.config.timeout)
        )
      ]);
      
      result = {
        name,
        status: 'pass',
        message: 'Test passed successfully',
        duration: Date.now() - startTime,
      };
    } catch (error) {
      result = {
        name,
        status: 'fail',
        message: 'Test failed',
        duration: Date.now() - startTime,
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }

    this.results.push(result);
    return result;
  }

  getResults(): TestResults {
    const passed = this.results.filter(r => r.status === 'pass').length;
    const failed = this.results.filter(r => r.status === 'fail').length;
    
    return {
      passed,
      failed,
      total: this.results.length,
      details: this.results,
    };
  }

  reset(): void {
    this.results = [];
  }
}

// API endpoint tests
export const runAPITests = async (config: TestConfig): Promise<TestResults> => {
  const tester = new IntegrationTester(config);
  const testClient = new TranzaAPIClient(config.baseURL);
  
  if (config.testToken) {
    testClient.setAuthToken(config.testToken);
  }

  // Test 1: Health Check
  await tester.runTest('API Health Check', async () => {
    const response = await fetch(`${config.baseURL}/health`);
    if (!response.ok) {
      throw new Error(`Health check failed: ${response.status}`);
    }
  });

  // Test 2: Authentication
  await tester.runTest('Authentication Test', async () => {
    if (!config.testToken) {
      throw new Error('No test token provided');
    }
    
    const response = await testClient.getUserProfile();
    if (!response.success) {
      throw new Error(`Auth failed: ${response.error}`);
    }
  });

  // Test 3: Wallet Balance
  await tester.runTest('Wallet Balance API', async () => {
    const response = await testClient.getWalletBalance();
    if (!response.success) {
      throw new Error(`Wallet balance failed: ${response.error}`);
    }
  });

  // Test 4: Transfer Validation
  await tester.runTest('Transfer Validation API', async () => {
    const response = await testClient.validateTransfer({
      amount: '100',
      recipient_type: 'upi',
      recipient_value: 'test@paytm',
    });
    
    // We expect this to work (validation logic handles invalid recipients)
    if (!response.success) {
      throw new Error(`Transfer validation failed: ${response.error}`);
    }
  });

  // Test 5: Transfer Fees
  await tester.runTest('Transfer Fees API', async () => {
    const response = await testClient.getTransferFees('100', 'upi');
    if (!response.success) {
      throw new Error(`Transfer fees failed: ${response.error}`);
    }
  });

  // Test 6: Transfer History
  await tester.runTest('Transfer History API', async () => {
    const response = await testClient.getTransferHistory(1, 5);
    if (!response.success) {
      throw new Error(`Transfer history failed: ${response.error}`);
    }
  });

  // Test 7: API Key Management
  await tester.runTest('API Key Management', async () => {
    const response = await testClient.getAPIKeys();
    if (!response.success) {
      throw new Error(`API key management failed: ${response.error}`);
    }
  });

  return tester.getResults();
};

// Form validation tests
export const runFormValidationTests = (): TestResults => {
  const tester = new IntegrationTester({ baseURL: '', timeout: 1000 });

  // Import validation functions
  const { validateAmount, validateUPI, validatePhone } = require('../services/api');

  // Test amount validation
  tester.runTest('Amount Validation - Valid', async () => {
    if (!validateAmount('100')) throw new Error('Valid amount rejected');
  });

  tester.runTest('Amount Validation - Invalid', async () => {
    if (validateAmount('0')) throw new Error('Zero amount accepted');
    if (validateAmount('-100')) throw new Error('Negative amount accepted');
    if (validateAmount('1000000')) throw new Error('Excessive amount accepted');
  });

  // Test UPI validation
  tester.runTest('UPI Validation - Valid', async () => {
    if (!validateUPI('user@paytm')) throw new Error('Valid UPI rejected');
    if (!validateUPI('test123@phonepe')) throw new Error('Valid UPI rejected');
  });

  tester.runTest('UPI Validation - Invalid', async () => {
    if (validateUPI('invalid-upi')) throw new Error('Invalid UPI accepted');
    if (validateUPI('user@')) throw new Error('Incomplete UPI accepted');
    if (validateUPI('@paytm')) throw new Error('Incomplete UPI accepted');
  });

  // Test phone validation
  tester.runTest('Phone Validation - Valid', async () => {
    if (!validatePhone('9876543210')) throw new Error('Valid phone rejected');
    if (!validatePhone('8765432109')) throw new Error('Valid phone rejected');
  });

  tester.runTest('Phone Validation - Invalid', async () => {
    if (validatePhone('1234567890')) throw new Error('Invalid phone accepted');
    if (validatePhone('98765432')) throw new Error('Short phone accepted');
    if (validatePhone('98765432101')) throw new Error('Long phone accepted');
  });

  return tester.getResults();
};

// End-to-end transfer test (simulation)
export const runTransferFlowTest = async (config: TestConfig): Promise<TestResults> => {
  const tester = new IntegrationTester(config);
  const testClient = new TranzaAPIClient(config.baseURL);
  
  if (config.testToken) {
    testClient.setAuthToken(config.testToken);
  }

  let transferId: string | undefined;

  // Test complete transfer flow
  await tester.runTest('Complete Transfer Flow', async () => {
    // Step 1: Validate transfer
    const validation = await testClient.validateTransfer({
      amount: '1', // Minimal amount for testing
      recipient_type: 'upi',
      recipient_value: 'test@paytm',
    });

    if (!validation.success) {
      throw new Error(`Validation failed: ${validation.error}`);
    }

    // Step 2: Create transfer (this might fail due to insufficient balance, which is expected)
    const creation = await testClient.createTransfer({
      amount: '1',
      recipient_type: 'upi',
      recipient_value: 'test@paytm',
      description: 'Test transfer',
    });

    // We don't expect this to succeed in a test environment
    // but we test that the API responds properly
    if (creation.success && creation.data) {
      transferId = creation.data.transfer_id;
    }
  });

  // Test transfer status if transfer was created
  if (transferId) {
    await tester.runTest('Transfer Status Check', async () => {
      const status = await testClient.getTransferStatus(transferId!);
      if (!status.success) {
        throw new Error(`Status check failed: ${status.error}`);
      }
    });
  }

  return tester.getResults();
};

// Performance tests
export const runPerformanceTests = async (config: TestConfig): Promise<TestResults> => {
  const tester = new IntegrationTester({ ...config, timeout: 10000 });
  const testClient = new TranzaAPIClient(config.baseURL);
  
  if (config.testToken) {
    testClient.setAuthToken(config.testToken);
  }

  // Test API response times
  await tester.runTest('API Response Time', async () => {
    const startTime = Date.now();
    const response = await testClient.getUserProfile();
    const duration = Date.now() - startTime;
    
    if (!response.success) {
      throw new Error(`API call failed: ${response.error}`);
    }
    
    if (duration > 5000) {
      throw new Error(`API response too slow: ${duration}ms`);
    }
  });

  // Test concurrent requests
  await tester.runTest('Concurrent Requests', async () => {
    const promises = Array(5).fill(null).map(() => testClient.getWalletBalance());
    const results = await Promise.all(promises);
    
    const failures = results.filter(r => !r.success);
    if (failures.length > 0) {
      throw new Error(`${failures.length} concurrent requests failed`);
    }
  });

  return tester.getResults();
};

// Utility function to run all tests
export const runAllTests = async (config: TestConfig): Promise<{
  api: TestResults;
  validation: TestResults;
  transfer: TestResults;
  performance: TestResults;
  overall: {
    totalPassed: number;
    totalFailed: number;
    totalTests: number;
    successRate: number;
  };
}> => {
  console.log('🧪 Running comprehensive test suite...');
  
  const api = await runAPITests(config);
  const validation = runFormValidationTests();
  const transfer = await runTransferFlowTest(config);
  const performance = await runPerformanceTests(config);

  const totalPassed = api.passed + validation.passed + transfer.passed + performance.passed;
  const totalFailed = api.failed + validation.failed + transfer.failed + performance.failed;
  const totalTests = totalPassed + totalFailed;
  const successRate = totalTests > 0 ? (totalPassed / totalTests) * 100 : 0;

  return {
    api,
    validation,
    transfer,
    performance,
    overall: {
      totalPassed,
      totalFailed,
      totalTests,
      successRate,
    },
  };
};
