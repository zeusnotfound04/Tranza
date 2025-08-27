'use client';

import React, { useState, useEffect } from 'react';
import { useRouter, useParams } from 'next/navigation';
import {
  ColumnDef,
  flexRender,
  getCoreRowModel,
  useReactTable,
  getSortedRowModel,
  SortingState,
  getPaginationRowModel,
} from '@tanstack/react-table';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  LineChart,
  Line,
  PieChart,
  Pie,
  Cell,
  Area,
  AreaChart,
} from 'recharts';
import { ArrowLeft, Clock, DollarSign, Activity, TrendingUp, AlertCircle, CheckCircle, XCircle } from 'lucide-react';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@tranza/ui/components/ui/card-ui';
import { Badge } from '@tranza/ui/components/ui/badge';
import { Button } from '@/components/ui/Button';
import { Alert, AlertDescription } from '@tranza/ui/components/ui/alert';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/Table';

import { APIKeyService } from '@/lib/services';
import { UsageStatsResponse, DetailedUsageResponse, APIUsageLog, TimeSeriesData, CommandUsage } from '@/types/api';
import { aeonikPro } from '@/lib/fonts';

const CHART_COLORS = ['#3b82f6', '#06b6d4', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899'];

export default function APIKeyUsagePage() {
  const router = useRouter();
  const params = useParams();
  const keyId = parseInt(params.id as string);

  const [usageStats, setUsageStats] = useState<UsageStatsResponse | null>(null);
  const [usageLogs, setUsageLogs] = useState<APIUsageLog[]>([]);
  const [timeSeriesData, setTimeSeriesData] = useState<TimeSeriesData[]>([]);
  const [commandData, setCommandData] = useState<CommandUsage[]>([]);
  
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [selectedTimeRange, setSelectedTimeRange] = useState(30);
  const [currentPage, setCurrentPage] = useState(0);
  const [sorting, setSorting] = useState<SortingState>([]);

  useEffect(() => {
    if (keyId) {
      loadUsageData();
    }
  }, [keyId, selectedTimeRange]);

  const loadUsageData = async () => {
    try {
      setLoading(true);
      setError('');

      // Load all data in parallel
      const [statsRes, logsRes, timeSeriesRes, commandRes] = await Promise.all([
        APIKeyService.getUsageStats(keyId, selectedTimeRange),
        APIKeyService.getUsageLogs(keyId, currentPage, 50),
        APIKeyService.getTimeSeriesData(keyId, selectedTimeRange),
        APIKeyService.getCommandData(keyId, selectedTimeRange),
      ]);

      if (statsRes.data) setUsageStats(statsRes.data);
      if (logsRes.data) setUsageLogs(logsRes.data.logs);
      if (timeSeriesRes.data) setTimeSeriesData(timeSeriesRes.data);
      if (commandRes.data) setCommandData(commandRes.data);

    } catch (err: any) {
      setError(err.message || 'Failed to load usage data');
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (amount: number, currency = 'INR') => {
    return new Intl.NumberFormat('en-IN', {
      style: 'currency',
      currency,
    }).format(amount);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-IN', {
      month: 'short',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const getStatusIcon = (statusCode: number) => {
    if (statusCode >= 200 && statusCode < 300) {
      return <CheckCircle className="w-4 h-4 text-green-500" />;
    } else if (statusCode >= 400) {
      return <XCircle className="w-4 h-4 text-red-500" />;
    } else {
      return <AlertCircle className="w-4 h-4 text-yellow-500" />;
    }
  };

  const getStatusBadgeVariant = (statusCode: number) => {
    if (statusCode >= 200 && statusCode < 300) return 'default';
    if (statusCode >= 400) return 'destructive';
    return 'secondary';
  };

  // Table columns definition
  const columns: ColumnDef<APIUsageLog>[] = [
    {
      accessorKey: 'timestamp',
      header: 'Time',
      cell: ({ row }) => (
        <div className="text-sm text-gray-300">
          {formatDate(row.getValue('timestamp'))}
        </div>
      ),
    },
    {
      accessorKey: 'command',
      header: 'Command',
      cell: ({ row }) => (
        <div className="font-mono text-sm bg-[#1a1a1a] px-2 py-1 rounded border border-gray-700">
          {row.getValue('command') || row.original.endpoint}
        </div>
      ),
    },
    {
      accessorKey: 'method',
      header: 'Method',
      cell: ({ row }) => (
        <Badge variant="outline" className="text-xs">
          {row.getValue('method')}
        </Badge>
      ),
    },
    {
      accessorKey: 'status_code',
      header: 'Status',
      cell: ({ row }) => {
        const statusCode = row.getValue('status_code') as number;
        return (
          <div className="flex items-center space-x-2">
            {getStatusIcon(statusCode)}
            <Badge variant={getStatusBadgeVariant(statusCode)} className="text-xs">
              {statusCode}
            </Badge>
          </div>
        );
      },
    },
    {
      accessorKey: 'amount_spent',
      header: 'Amount',
      cell: ({ row }) => {
        const amount = row.getValue('amount_spent') as number;
        if (!amount) return <span className="text-gray-500">-</span>;
        return (
          <span className={amount > 0 ? 'text-red-400' : 'text-gray-400'}>
            {formatCurrency(amount)}
          </span>
        );
      },
    },
    {
      accessorKey: 'response_time',
      header: 'Response Time',
      cell: ({ row }) => (
        <span className="text-sm text-gray-300">
          {row.getValue('response_time')}ms
        </span>
      ),
    },
  ];

  const table = useReactTable({
    data: usageLogs,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    onSortingChange: setSorting,
    state: {
      sorting,
    },
  });

  if (loading) {
    return (
      <div className={`min-h-screen bg-[#121212] p-6 ${aeonikPro.className}`}>
        <div className="flex items-center justify-center h-96">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={`min-h-screen bg-[#121212] p-6 ${aeonikPro.className}`}>
        <div className="mb-6">
          <Button onClick={() => router.back()} variant="ghost" className="text-gray-300 hover:text-white">
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back
          </Button>
        </div>
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      </div>
    );
  }

  return (
    <div className={`min-h-screen bg-[#121212] p-6 ${aeonikPro.className}`}>
      {/* Header */}
      <div className="mb-6">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <Button onClick={() => router.back()} variant="ghost" className="text-gray-300 hover:text-white">
              <ArrowLeft className="w-4 h-4 mr-2" />
              Back
            </Button>
            <div>
              <h1 className="text-3xl font-bold text-white">API Key Usage</h1>
              <p className="text-gray-400">Detailed analytics and usage logs for API key #{keyId}</p>
            </div>
          </div>
          
          {/* Time Range Selector */}
          <div className="flex items-center space-x-2">
            {[7, 30, 90].map((days) => (
              <Button
                key={days}
                variant={selectedTimeRange === days ? "default" : "outline"}
                size="sm"
                onClick={() => setSelectedTimeRange(days)}
                className={selectedTimeRange === days ? "bg-blue-600 hover:bg-blue-700" : ""}
              >
                {days}d
              </Button>
            ))}
          </div>
        </div>
      </div>

      {/* Statistics Cards */}
      {usageStats && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <Card className="bg-[#1a1a1a] border-gray-700">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-gray-300">Total Requests</CardTitle>
              <Activity className="h-4 w-4 text-blue-400" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-white">{usageStats.total_requests.toLocaleString()}</div>
              <p className="text-xs text-gray-400">
                Success rate: {(usageStats.success_rate * 100).toFixed(1)}%
              </p>
            </CardContent>
          </Card>

          <Card className="bg-[#1a1a1a] border-gray-700">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-gray-300">Total Spent</CardTitle>
              <DollarSign className="h-4 w-4 text-green-400" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-white">
                {formatCurrency(usageStats.total_amount_spent)}
              </div>
              <p className="text-xs text-gray-400">
                Limit: {formatCurrency(usageStats.spending_limit)}
              </p>
            </CardContent>
          </Card>

          <Card className="bg-[#1a1a1a] border-gray-700">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-gray-300">Remaining Limit</CardTitle>
              <TrendingUp className="h-4 w-4 text-purple-400" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-white">
                {formatCurrency(usageStats.remaining_limit)}
              </div>
              <p className="text-xs text-gray-400">
                {((usageStats.remaining_limit / usageStats.spending_limit) * 100).toFixed(1)}% remaining
              </p>
            </CardContent>
          </Card>

          <Card className="bg-[#1a1a1a] border-gray-700">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-gray-300">Avg Response Time</CardTitle>
              <Clock className="h-4 w-4 text-orange-400" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-white">{usageStats.avg_response_time.toFixed(0)}ms</div>
              <p className="text-xs text-gray-400">
                Last used: {usageStats.last_used_at ? formatDate(usageStats.last_used_at) : 'Never'}
              </p>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Charts Section */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        {/* Usage Over Time */}
        <Card className="bg-[#1a1a1a] border-gray-700">
          <CardHeader>
            <CardTitle className="text-white">Usage Over Time</CardTitle>
            <CardDescription className="text-gray-400">Requests and spending trends</CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <AreaChart data={timeSeriesData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                <XAxis 
                  dataKey="date" 
                  stroke="#9ca3af"
                  fontSize={12}
                  tickFormatter={(value) => new Date(value).toLocaleDateString('en-IN', { month: 'short', day: 'numeric' })}
                />
                <YAxis stroke="#9ca3af" fontSize={12} />
                <Tooltip 
                  contentStyle={{ 
                    backgroundColor: '#1f2937', 
                    border: '1px solid #374151',
                    borderRadius: '8px',
                    color: '#f3f4f6'
                  }}
                />
                <Area
                  type="monotone"
                  dataKey="requests"
                  stackId="1"
                  stroke="#3b82f6"
                  fill="#3b82f6"
                  fillOpacity={0.3}
                />
              </AreaChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* Command Usage Distribution */}
        <Card className="bg-[#1a1a1a] border-gray-700">
          <CardHeader>
            <CardTitle className="text-white">Command Usage</CardTitle>
            <CardDescription className="text-gray-400">Distribution by command type</CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={commandData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ command, count }) => `${command}: ${count}`}
                  outerRadius={80}
                  fill="#8884d8"
                  dataKey="count"
                >
                  {commandData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={CHART_COLORS[index % CHART_COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip 
                  contentStyle={{ 
                    backgroundColor: '#1f2937', 
                    border: '1px solid #374151',
                    borderRadius: '8px',
                    color: '#f3f4f6'
                  }}
                />
              </PieChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </div>

      {/* Spending by Command */}
      <Card className="bg-[#1a1a1a] border-gray-700 mb-8">
        <CardHeader>
          <CardTitle className="text-white">Spending by Command</CardTitle>
          <CardDescription className="text-gray-400">Amount spent per command type</CardDescription>
        </CardHeader>
        <CardContent>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={commandData}>
              <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
              <XAxis dataKey="command" stroke="#9ca3af" fontSize={12} />
              <YAxis stroke="#9ca3af" fontSize={12} />
              <Tooltip 
                contentStyle={{ 
                  backgroundColor: '#1f2937', 
                  border: '1px solid #374151',
                  borderRadius: '8px',
                  color: '#f3f4f6'
                }}
                formatter={(value: any) => [formatCurrency(value), 'Amount']}
              />
              <Bar dataKey="total_amount" fill="#10b981" />
            </BarChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>

      {/* Usage Logs Table */}
      <Card className="bg-[#1a1a1a] border-gray-700">
        <CardHeader>
          <CardTitle className="text-white">Recent Activity</CardTitle>
          <CardDescription className="text-gray-400">Detailed logs of API usage</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="rounded-md border border-gray-700">
            <Table>
              <TableHeader>
                {table.getHeaderGroups().map((headerGroup) => (
                  <TableRow key={headerGroup.id} className="border-gray-700 hover:bg-[#1a1a1a]">
                    {headerGroup.headers.map((header) => (
                      <TableHead key={header.id} className="text-gray-300">
                        {header.isPlaceholder
                          ? null
                          : flexRender(
                              header.column.columnDef.header,
                              header.getContext()
                            )}
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
                      className="border-gray-700 hover:bg-[#1f1f1f]"
                      data-state={row.getIsSelected() && "selected"}
                    >
                      {row.getVisibleCells().map((cell) => (
                        <TableCell key={cell.id} className="text-gray-200">
                          {flexRender(cell.column.columnDef.cell, cell.getContext())}
                        </TableCell>
                      ))}
                    </TableRow>
                  ))
                ) : (
                  <TableRow>
                    <TableCell colSpan={columns.length} className="h-24 text-center text-gray-400">
                      No usage data found.
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
                usageLogs.length
              )}{' '}
              of {usageLogs.length} entries
            </div>
            <div className="space-x-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => table.previousPage()}
                disabled={!table.getCanPreviousPage()}
              >
                Previous
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => table.nextPage()}
                disabled={!table.getCanNextPage()}
              >
                Next
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
