import React, { useState, useEffect } from 'react';
import {
  Container,
  Paper,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Box,
  Alert,
  Button,
  CircularProgress,
  Chip,
  Divider,
  Card,
  CardContent,
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { getHistory } from '../services/api';

// Updated interface to match backend model
interface Transaction {
  id: number;
  account_id: number;
  type: string; // "EARN" or "REDEEM"
  points: number;
  transaction_at: string;
  created_at: string;
}

const TransactionHistory: React.FC = () => {
  const navigate = useNavigate();
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [error, setError] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchHistory = async () => {
      try {
        console.log('Fetching transaction history...');
        const response = await getHistory();
        
        if (response && response.transactions && Array.isArray(response.transactions)) {
          setTransactions(response.transactions);
        } else {
          console.error('Invalid transactions data format:', response);
          setError('Invalid transaction data format');
        }
      } catch (err) {
        console.error('Error fetching history:', err);
        setError('Failed to fetch transaction history');
      } finally {
        setLoading(false);
      }
    };

    fetchHistory();
  }, []);

  const formatDate = (dateString: string) => {
    if (!dateString) return 'No date available';
    
    try {
      const date = new Date(dateString);
      if (isNaN(date.getTime())) {
        return 'Invalid date format';
      }
      
      // Format as: "Apr 18, 2025 at 3:29 PM"
      return new Intl.DateTimeFormat('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
        hour: 'numeric',
        minute: 'numeric',
        hour12: true
      }).format(date);
    } catch (e) {
      console.error('Date parsing error:', e, 'for value:', dateString);
      return 'Date error';
    }
  };

  const getTotalPoints = () => {
    if (!transactions.length) return 0;
    
    return transactions.reduce((total, transaction) => {
      if (transaction.type === 'EARN') {
        return total + transaction.points;
      } else {
        return total - transaction.points;
      }
    }, 0);
  };

  const getEarnedPoints = () => {
    if (!transactions.length) return 0;
    
    return transactions
      .filter(t => t.type === 'EARN')
      .reduce((total, t) => total + t.points, 0);
  };

  const getRedeemedPoints = () => {
    if (!transactions.length) return 0;
    
    return transactions
      .filter(t => t.type === 'REDEEM')
      .reduce((total, t) => total + t.points, 0);
  };

  return (
    <Container maxWidth="md">
      <Box sx={{ mt: 4 }}>
        <Paper sx={{ p: 4 }}>
          <Typography variant="h4" gutterBottom>
            Transaction History
          </Typography>

          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}

          {loading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', my: 4 }}>
              <CircularProgress />
            </Box>
          ) : transactions.length === 0 ? (
            <Alert severity="info">
              No transactions found. Earn or redeem points to see your history.
            </Alert>
          ) : (
            <>
              {/* Summary Cards */}
              <Box 
                sx={{ 
                  display: 'flex', 
                  gap: 2, 
                  mb: 4, 
                  flexDirection: { xs: 'column', sm: 'row' } 
                }}
              >
                <Card sx={{ flex: 1, bgcolor: '#f8f8f8' }}>
                  <CardContent>
                    <Typography color="textSecondary" gutterBottom>
                      Total Earned
                    </Typography>
                    <Typography variant="h5" component="div" sx={{ color: 'green' }}>
                      +{getEarnedPoints()} points
                    </Typography>
                  </CardContent>
                </Card>
                
                <Card sx={{ flex: 1, bgcolor: '#f8f8f8' }}>
                  <CardContent>
                    <Typography color="textSecondary" gutterBottom>
                      Total Redeemed
                    </Typography>
                    <Typography variant="h5" component="div" sx={{ color: 'red' }}>
                      -{getRedeemedPoints()} points
                    </Typography>
                  </CardContent>
                </Card>
                
                <Card sx={{ flex: 1, bgcolor: '#e8f4ff' }}>
                  <CardContent>
                    <Typography color="textSecondary" gutterBottom>
                      Net Balance from Transactions
                    </Typography>
                    <Typography 
                      variant="h5" 
                      component="div" 
                      sx={{ 
                        color: getTotalPoints() >= 0 ? 'green' : 'red' 
                      }}
                    >
                      {getTotalPoints() >= 0 ? '+' : ''}{getTotalPoints()} points
                    </Typography>
                  </CardContent>
                </Card>
              </Box>

              <Divider sx={{ my: 3 }} />
              
              <Typography variant="h6" gutterBottom>
                Detailed Transactions
              </Typography>
              
              <TableContainer>
                <Table>
                  <TableHead sx={{ bgcolor: '#f5f5f5' }}>
                    <TableRow>
                      <TableCell>Date & Time</TableCell>
                      <TableCell>Transaction Type</TableCell>
                      <TableCell align="right">Points</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {transactions.map((transaction) => (
                      <TableRow key={transaction.id} 
                        sx={{ 
                          '&:hover': { 
                            bgcolor: '#f9f9f9' 
                          },
                          borderLeft: '4px solid',
                          borderLeftColor: transaction.type === 'EARN' ? 'success.light' : 'error.light',
                        }}
                      >
                        <TableCell>
                          {formatDate(transaction.transaction_at || transaction.created_at)}
                        </TableCell>
                        <TableCell>
                          <Chip
                            label={transaction.type === 'EARN' ? 'Earned' : 'Redeemed'}
                            color={transaction.type === 'EARN' ? 'success' : 'error'}
                            variant="outlined"
                            size="small"
                          />
                        </TableCell>
                        <TableCell align="right">
                          <Typography
                            sx={{ 
                              fontWeight: 'bold',
                              color: transaction.type === 'EARN' ? 'green' : 'red' 
                            }}
                          >
                            {transaction.type === 'EARN' ? '+' : '-'}
                            {transaction.points}
                          </Typography>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            </>
          )}

          <Box sx={{ mt: 4 }}>
            <Button
              variant="contained"
              onClick={() => navigate('/')}
              fullWidth
            >
              Back to Dashboard
            </Button>
          </Box>
        </Paper>
      </Box>
    </Container>
  );
};

export default TransactionHistory; 