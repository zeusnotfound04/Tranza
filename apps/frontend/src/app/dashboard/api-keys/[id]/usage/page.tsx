'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import {
  LineChart,
  Line,
  AreaChart,
  Area,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import {
  useReactTable,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  createColumnHelper,
  flexRender,
  ColumnDef,
  SortingState,
} from '@tanstack/react-table';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@tranza/ui/components/ui/card-ui';
import { Button } from '@/components/ui/Button';
import { Badge } from '@tranza/ui/components/ui/badge';
import { Alert, AlertDescription } from '@tranza/ui/components/ui/alert';
import { Input } from '@tranza/ui/components/ui/input';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/Table';
import {
  ArrowLeft,
  Activity,
  DollarSign,
  Clock,
  TrendingUp,
  ChevronLeft,
  ChevronRight,
  Search,
  Filter,
  Download,
  RefreshCw,
} from 'lucide-react';
import { aeonikPro } from '@/lib/fonts';
import { apiClient } from '@/services/api';
import { APIUsageLog, UsageStatsResponse, CommandUsage, TimeSeriesData } from '@/types/api';

const columnHelper = createColumnHelper<APIUsageLog>();

export default function APIKeyUsagePage() {
  const params = useParams();
  const router = useRouter();
  const keyId = params.id as string;

  // State management
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [refreshing, setRefreshing] = useState(false);
  
  // Data states
  const [usageStats, setUsageStats] = useState<UsageStatsResponse | null>(null);
  const [usageLogs, setUsageLogs] = useState<APIUsageLog[]>([]);
  const [commandData, setCommandData] = useState<CommandUsage[]>([]);
  const [timeSeriesData, setTimeSeriesData] = useState<TimeSeriesData[]>([]);
  
  // Table states
  const [globalFilter, setGlobalFilter] = useState('');
  const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 10 });
  const [sorting, setSorting] = useState<SortingState>([]);
  
  // Filter states
  const [statusFilter, setStatusFilter] = useState('all');
  const [commandFilter, setCommandFilter] = useState('all');
  const [timePeriod, setTimePeriod] = useState('7d');

  // Table columns definition
  const columns: ColumnDef<APIUsageLog, any>[] = [
    {
      accessorKey: 'timestamp',
      header: 'Time',
      cell: ({ row }) => (
        <div className="flex flex-col">
          <span className="text-sm font-medium text-white">
            {new Date(row.original.timestamp).toLocaleDateString()}
          </span>
          <span className="text-xs text-gray-400">
            {new Date(row.original.timestamp).toLocaleTimeString()}
          </span>
        </div>
      ),
    },
    {
      accessorKey: 'command',
      header: 'Command',
      cell: ({ row }) => (
        <Badge 
          variant="outline" 
          className="bg-[#1a1a1a] border-gray-600 text-blue-400"
        >
          {row.original.command || row.original.endpoint}
        </Badge>
      ),
    },
    {
      accessorKey: 'method',
      header: 'Method',
      cell: ({ row }) => (
        <Badge 
          variant={row.original.method === 'GET' ? 'default' : 'secondary'}
          className="text-xs"
        >
          {row.original.method}
        </Badge>
      ),
    },
    {
      accessorKey: 'status_code',
      header: 'Status',
      cell: ({ row }) => {
        const status = row.original.status_code;
        const variant = status >= 200 && status < 300 ? 'default' : status >= 400 ? 'destructive' : 'secondary';
        return (
          <Badge variant={variant} className="text-xs">
            {status}
          </Badge>
        );
      },
    },
    {
      accessorKey: 'amount_spent',
      header: 'Amount',
      cell: ({ row }) => (
        <div className="text-right">
          {row.original.amount_spent ? (
            <span className="text-sm font-medium text-red-400">
              -₹{row.original.amount_spent.toFixed(2)}
            </span>
          ) : (
            <span className="text-xs text-gray-500">-</span>
          )}
        </div>
      ),
    },
    {
      accessorKey: 'response_time',
      header: 'Response Time',
      cell: ({ row }) => (
        <span className="text-xs text-gray-400">
          {row.original.response_time}ms
        </span>
      ),
    },
    {
      accessorKey: 'ip_address',
      header: 'IP Address',
      cell: ({ row }) => (
        <span className="text-xs text-gray-400 font-mono">
          {row.original.ip_address}
        </span>
      ),
    },
  ];

  const table = useReactTable({
    data: usageLogs,
    columns,
    state: {
      globalFilter,
      pagination,
      sorting,
    },
    onGlobalFilterChange: setGlobalFilter,
    onPaginationChange: setPagination,
    onSortingChange: setSorting,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
  });

  // Load data functions
  const loadUsageStats = async () => {
    try {
      const response = await apiClient.getDetailedUsageStats(keyId);
      if (response.data) {
        setUsageStats(response.data);
      }
    } catch (err) {
      console.error('Error loading usage stats:', err);
      setError('Failed to load usage statistics');
    }
  };

  const loadUsageLogs = async () => {
    try {
      const response = await apiClient.getUsageLogs(keyId, 0, 100);
      if (response.data) {
        setUsageLogs(response.data.logs || []);
      }
    } catch (err) {
      console.error('Error loading usage logs:', err);
      setError('Failed to load usage logs');
    }
  };

  const loadCommandData = async () => {
    try {
      const response = await apiClient.getCommandData(keyId);
      if (response.data) {
        setCommandData(response.data.commands || []);
      }
    } catch (err) {
      console.error('Error loading command data:', err);
    }
  };

  const loadTimeSeriesData = async () => {
    try {
      const response = await apiClient.getTimeSeriesData(keyId, timePeriod);
      if (response.data) {
        setTimeSeriesData(response.data.time_series || []);
      }
    } catch (err) {
      console.error('Error loading time series data:', err);
    }
  };

  const loadAllData = async () => {
    setLoading(true);
    setError('');
    
    try {
      await Promise.all([
        loadUsageStats(),
        loadUsageLogs(),
        loadCommandData(),
        loadTimeSeriesData(),
      ]);
    } catch (err) {
      console.error('Error loading data:', err);
      setError('Failed to load API key usage data');
    } finally {
      setLoading(false);
    }
  };

  const refreshData = async () => {
    setRefreshing(true);
    await loadAllData();
    setRefreshing(false);
  };

  useEffect(() => {
    if (keyId) {
      loadAllData();
    }
  }, [keyId, timePeriod]);

  // Chart colors for dark theme
  const chartColors = {
    primary: '#3b82f6',
    secondary: '#10b981',
    accent: '#f59e0b',
    danger: '#ef4444',
    success: '#22c55e',
    warning: '#f59e0b',
    info: '#06b6d4',
  };

  const pieColors = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#06b6d4'];

  // Format currency
  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-IN', {
      style: 'currency',
      currency: 'INR',
    }).format(amount);
  };

  // Calculate spending percentage
  const getSpendingPercentage = () => {
    if (!usageStats) return 0;
    return (usageStats.total_amount_spent / usageStats.spending_limit) * 100;
  };

  if (loading) {
    return (
      <div className={`min-h-screen bg-[#121212] p-6 ${aeonikPro.className}`}>
        <div className="flex items-center justify-center py-20">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
        </div>
      </div>
    );
  }

  if (error && !usageStats) {
    return (
      <div className={`min-h-screen bg-[#121212] p-6 ${aeonikPro.className}`}>
        <Alert variant="destructive" className="max-w-2xl mx-auto mt-20">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
        <div className="flex justify-center mt-4">
          <Button onClick={() => router.back()} variant="outline">
            <ArrowLeft className="w-4 h-4 mr-2" />
            Go Back
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className={`min-h-screen bg-[#121212] p-6 ${aeonikPro.className}`}>
      {/* Header */}
      <div className="flex items-center justify-between mb-8">
        <div className="flex items-center space-x-4">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => router.push('/dashboard/api-keys')}
            className="text-gray-400 hover:text-white"
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back to API Keys
          </Button>
          <div>
            <h1 className="text-3xl font-bold text-white">API Key Usage</h1>
            <p className="text-gray-400 mt-1">
              Detailed usage analytics and logs for API Key #{keyId}
            </p>
          </div>
        </div>
        <div className="flex items-center space-x-2">
          <select 
            value={timePeriod} 
            onChange={(e) => setTimePeriod(e.target.value)}
            className="w-32 bg-[#1f1f1f] border border-gray-700 text-white px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="24h">24 Hours</option>
            <option value="7d">7 Days</option>
            <option value="30d">30 Days</option>
            <option value="90d">90 Days</option>
          </select>
          <Button
            variant="outline"
            size="sm"
            onClick={refreshData}
            disabled={refreshing}
            className="border-gray-700 text-gray-300 hover:text-white"
          >
            <RefreshCw className={`w-4 h-4 mr-2 ${refreshing ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      {/* Overview Cards */}
      {usageStats && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <Card className="bg-[#1f1f1f] border-gray-700">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-gray-400">Total Spent</CardTitle>
              <DollarSign className="h-4 w-4 text-red-400" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-white">
                {formatCurrency(usageStats.total_amount_spent)}
              </div>
              <div className="text-xs text-gray-500 mt-1">
                of {formatCurrency(usageStats.spending_limit)} limit
              </div>
              <div className="w-full bg-gray-700 rounded-full h-2 mt-2">
                <div
                  className="bg-red-500 h-2 rounded-full transition-all duration-300"
                  style={{ width: `${Math.min(getSpendingPercentage(), 100)}%` }}
                />
              </div>
            </CardContent>
          </Card>

          <Card className="bg-[#1f1f1f] border-gray-700">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-gray-400">Total Requests</CardTitle>
              <Activity className="h-4 w-4 text-blue-400" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-white">
                {usageStats.total_requests.toLocaleString()}
              </div>
              <div className="text-xs text-gray-500 mt-1">
                Success Rate: {(usageStats.success_rate * 100).toFixed(1)}%
              </div>
            </CardContent>
          </Card>

          <Card className="bg-[#1f1f1f] border-gray-700">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-gray-400">Avg Response Time</CardTitle>
              <Clock className="h-4 w-4 text-yellow-400" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-white">
                {usageStats.avg_response_time.toFixed(0)}ms
              </div>
              <div className="text-xs text-gray-500 mt-1">
                Last used: {usageStats.last_used_at ? new Date(usageStats.last_used_at).toLocaleDateString() : 'Never'}
              </div>
            </CardContent>
          </Card>

          <Card className="bg-[#1f1f1f] border-gray-700">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-gray-400">Remaining Limit</CardTitle>
              <TrendingUp className="h-4 w-4 text-green-400" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-white">
                {formatCurrency(usageStats.remaining_limit)}
              </div>
              <div className="text-xs text-gray-500 mt-1">
                Available for spending
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Charts Section */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        {/* Usage Trend Chart */}
        <Card className="bg-[#1f1f1f] border-gray-700">
          <CardHeader>
            <CardTitle className="text-white">Usage Trend</CardTitle>
            <CardDescription className="text-gray-400">
              API requests and spending over time
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <AreaChart data={timeSeriesData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                <XAxis 
                  dataKey="date" 
                  stroke="#9ca3af"
                  fontSize={12}
                  tickFormatter={(value) => new Date(value).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}
                />
                <YAxis stroke="#9ca3af" fontSize={12} />
                <Tooltip 
                  contentStyle={{ 
                    backgroundColor: '#1f1f1f', 
                    border: '1px solid #374151',
                    borderRadius: '8px',
                    color: '#ffffff'
                  }}
                />
                <Legend />
                <Area
                  type="monotone"
                  dataKey="requests"
                  stackId="1"
                  stroke={chartColors.primary}
                  fill={chartColors.primary}
                  fillOpacity={0.6}
                  name="Requests"
                />
                <Area
                  type="monotone"
                  dataKey="amount_spent"
                  stackId="2"
                  stroke={chartColors.danger}
                  fill={chartColors.danger}
                  fillOpacity={0.6}
                  name="Amount Spent (₹)"
                />
              </AreaChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* Command Usage Distribution */}
        <Card className="bg-[#1f1f1f] border-gray-700">
          <CardHeader>
            <CardTitle className="text-white">Command Usage</CardTitle>
            <CardDescription className="text-gray-400">
              Distribution of API commands used
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={commandData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                  outerRadius={80}
                  fill="#8884d8"
                  dataKey="count"
                  nameKey="command"
                >
                  {commandData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={pieColors[index % pieColors.length]} />
                  ))}
                </Pie>
                <Tooltip 
                  contentStyle={{ 
                    backgroundColor: '#1f1f1f', 
                    border: '1px solid #374151',
                    borderRadius: '8px',
                    color: '#ffffff'
                  }}
                />
              </PieChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </div>

      {/* Response Time Chart */}
      <Card className="bg-[#1f1f1f] border-gray-700 mb-8">
        <CardHeader>
          <CardTitle className="text-white">Performance Metrics</CardTitle>
          <CardDescription className="text-gray-400">
            Response time and success rate trends
          </CardDescription>
        </CardHeader>
        <CardContent>
          <ResponsiveContainer width="100%" height={300}>
            <LineChart data={timeSeriesData}>
              <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
              <XAxis 
                dataKey="date" 
                stroke="#9ca3af"
                fontSize={12}
                tickFormatter={(value) => new Date(value).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}
              />
              <YAxis yAxisId="left" stroke="#9ca3af" fontSize={12} />
              <YAxis yAxisId="right" orientation="right" stroke="#9ca3af" fontSize={12} />
              <Tooltip 
                contentStyle={{ 
                  backgroundColor: '#1f1f1f', 
                  border: '1px solid #374151',
                  borderRadius: '8px',
                  color: '#ffffff'
                }}
              />
              <Legend />
              <Line
                yAxisId="left"
                type="monotone"
                dataKey="avg_response_time"
                stroke={chartColors.warning}
                strokeWidth={2}
                dot={{ fill: chartColors.warning, strokeWidth: 2, r: 4 }}
                name="Avg Response Time (ms)"
              />
              <Line
                yAxisId="right"
                type="monotone"
                dataKey="success_rate"
                stroke={chartColors.success}
                strokeWidth={2}
                dot={{ fill: chartColors.success, strokeWidth: 2, r: 4 }}
                name="Success Rate (%)"
              />
            </LineChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>

      {/* Usage Logs Table */}
      <Card className="bg-[#1f1f1f] border-gray-700">
        <CardHeader>
          <CardTitle className="text-white">Request Logs</CardTitle>
          <CardDescription className="text-gray-400">
            Detailed logs of all API requests made with this key
          </CardDescription>
          
          {/* Table Controls */}
          <div className="flex flex-col sm:flex-row gap-4 mt-4">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
              <Input
                placeholder="Search logs..."
                value={globalFilter}
                onChange={(e) => setGlobalFilter(e.target.value)}
                className="pl-10 bg-[#2a2a2a] border-gray-600 text-white placeholder-gray-400"
              />
            </div>
            <select 
              value={statusFilter} 
              onChange={(e) => setStatusFilter(e.target.value)}
              className="w-40 bg-[#2a2a2a] border border-gray-600 text-white px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="all">All Status</option>
              <option value="success">Success (2xx)</option>
              <option value="error">Error (4xx/5xx)</option>
            </select>
            <select 
              value={commandFilter} 
              onChange={(e) => setCommandFilter(e.target.value)}
              className="w-40 bg-[#2a2a2a] border border-gray-600 text-white px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="all">All Commands</option>
              {commandData.map((cmd) => (
                <option key={cmd.command} value={cmd.command}>
                  {cmd.command}
                </option>
              ))}
            </select>
          </div>
        </CardHeader>
        <CardContent>
          {/* Table */}
          <div className="rounded-md border border-gray-700 overflow-hidden">
            <Table>
              <TableHeader>
                {table.getHeaderGroups().map((headerGroup) => (
                  <TableRow key={headerGroup.id} className="border-gray-700 hover:bg-[#2a2a2a]">
                    {headerGroup.headers.map((header) => (
                      <TableHead 
                        key={header.id}
                        className="text-gray-300 font-medium cursor-pointer hover:text-white"
                        onClick={header.column.getToggleSortingHandler()}
                      >
                        {header.isPlaceholder
                          ? null
                          : flexRender(header.column.columnDef.header, header.getContext())}
                      </TableHead>
                    ))}
                  </TableRow>
                ))}
              </TableHeader>
              <TableBody>
                {table.getRowModel().rows?.length ? (
                  table.getRowModel().rows.map((row) => (
                    <TableRow 
                      key={row.id}
                      className="border-gray-700 hover:bg-[#2a2a2a] transition-colors"
                    >
                      {row.getVisibleCells().map((cell) => (
                        <TableCell key={cell.id}>
                          {flexRender(cell.column.columnDef.cell, cell.getContext())}
                        </TableCell>
                      ))}
                    </TableRow>
                  ))
                ) : (
                  <TableRow>
                    <TableCell 
                      colSpan={columns.length} 
                      className="h-24 text-center text-gray-400"
                    >
                      No usage logs found.
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>

          {/* Pagination */}
          <div className="flex items-center justify-between space-x-2 py-4">
            <div className="text-sm text-gray-400">
              Showing {table.getState().pagination.pageIndex * table.getState().pagination.pageSize + 1} to{' '}
              {Math.min(
                (table.getState().pagination.pageIndex + 1) * table.getState().pagination.pageSize,
                table.getFilteredRowModel().rows.length
              )}{' '}
              of {table.getFilteredRowModel().rows.length} entries
            </div>
            <div className="flex items-center space-x-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => table.previousPage()}
                disabled={!table.getCanPreviousPage()}
                className="border-gray-600 text-gray-300 hover:text-white"
              >
                <ChevronLeft className="h-4 w-4" />
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => table.nextPage()}
                disabled={!table.getCanNextPage()}
                className="border-gray-600 text-gray-300 hover:text-white"
              >
                <ChevronRight className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
