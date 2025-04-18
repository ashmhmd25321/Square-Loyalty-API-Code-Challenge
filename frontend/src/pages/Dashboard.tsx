import React, { useState, useEffect } from 'react';
import {
  Container,
  Paper,
  Typography,
  Button,
  TextField,
  Box,
  Grid,
  Alert,
  Chip,
  Avatar,
  Divider,
  InputAdornment,
  Card,
  CardContent,
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { getBalance, earnPoints, redeemPoints } from '../services/api';
import { getUserInfo } from '../services/api';

const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const [balance, setBalance] = useState<number>(0);
  const [points, setPoints] = useState<string>('');
  const [error, setError] = useState<string>('');
  const [success, setSuccess] = useState<string>('');
  const [userEmail, setUserEmail] = useState<string>('');
  const [purchaseAmount, setPurchaseAmount] = useState<string>('');
  const [pointsToEarn, setPointsToEarn] = useState<number>(0);

  const POINTS_CONVERSION_RATE = 10 / 1000;

  useEffect(() => {
    const amount = parseFloat(purchaseAmount);
    if (!isNaN(amount) && amount > 0) {
      setPointsToEarn(Math.floor(amount * POINTS_CONVERSION_RATE));
    } else {
      setPointsToEarn(0);
    }
  }, [purchaseAmount]);

  useEffect(() => {
    const userInfo = getUserInfo();
    if (userInfo) {
      setUserEmail(userInfo.email);
    }

    const fetchBalance = async () => {
      try {
        const response = await getBalance();
        setBalance(response.balance);
      } catch (err) {
        setError('Failed to fetch balance');
      }
    };

    fetchBalance();
  }, []);

  const handlePurchase = async () => {
    try {
      const amount = parseFloat(purchaseAmount);
      if (isNaN(amount) || amount <= 0) {
        setError('Please enter a valid purchase amount');
        return;
      }

      const pointsToAdd = Math.floor(amount * POINTS_CONVERSION_RATE);
      
      await earnPoints(pointsToAdd);
      const response = await getBalance();
      setBalance(response.balance);
      setSuccess(`Purchase of LKR ${amount.toLocaleString()} successful! You earned ${pointsToAdd} points.`);
      setPurchaseAmount('');
      setError('');
    } catch (err) {
      setError('Failed to process purchase');
      setSuccess('');
    }
  };

  const handleRedeemPoints = async () => {
    try {
      const pointsNum = parseInt(points);
      if (isNaN(pointsNum) || pointsNum <= 0) {
        setError('Please enter a valid number of points');
        return;
      }

      if (pointsNum > balance) {
        setError('Not enough points to redeem');
        return;
      }

      await redeemPoints(pointsNum);
      const response = await getBalance();
      setBalance(response.balance);
      setSuccess('Points redeemed successfully!');
      setPoints('');
      setError('');
    } catch (err) {
      setError('Failed to redeem points');
      setSuccess('');
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user_email');
    navigate('/login');
  };

  return (
    <Container maxWidth="md">
      <Box sx={{ mt: 4 }}>
        <Paper sx={{ p: 4 }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
            <Typography variant="h4">
              Loyalty Program
            </Typography>
            <Chip
              avatar={<Avatar>{userEmail.charAt(0).toUpperCase()}</Avatar>}
              label={userEmail}
              variant="outlined"
              color="primary"
            />
          </Box>

          <Divider sx={{ mb: 3 }} />

          <Box sx={{ 
            mb: 4, 
            p: 3, 
            bgcolor: '#f5f9ff', 
            borderRadius: 2,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center' 
          }}>
            <Typography variant="h6" color="textSecondary" gutterBottom>
              Current Balance
            </Typography>
            <Typography variant="h3" sx={{ fontWeight: 'bold', color: '#2196f3' }}>
              {balance} points
            </Typography>
          </Box>

          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}

          {success && (
            <Alert severity="success" sx={{ mb: 2 }}>
              {success}
            </Alert>
          )}

          <Card sx={{ mb: 4, bgcolor: '#f9f9f9' }}>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Make a Purchase
              </Typography>
              <Grid container spacing={2} alignItems="center">
                <Grid item xs={12} sm={6}>
                  <TextField
                    fullWidth
                    label="Purchase Amount"
                    type="number"
                    InputProps={{
                      startAdornment: <InputAdornment position="start">LKR</InputAdornment>,
                    }}
                    value={purchaseAmount}
                    onChange={(e) => setPurchaseAmount(e.target.value)}
                    margin="normal"
                  />
                </Grid>
                <Grid item xs={12} sm={6}>
                  <Box sx={{ display: 'flex', flexDirection: 'column' }}>
                    <Typography variant="body1" color="textSecondary" gutterBottom>
                      Points you'll earn:
                    </Typography>
                    <Typography variant="h5" color="primary" fontWeight="bold">
                      {pointsToEarn} points
                    </Typography>
                    <Typography variant="caption" color="textSecondary">
                      (10 points per 1,000 LKR spent)
                    </Typography>
                  </Box>
                </Grid>
                <Grid item xs={12}>
                  <Button
                    variant="contained"
                    color="primary"
                    onClick={handlePurchase}
                    fullWidth
                    disabled={pointsToEarn <= 0}
                  >
                    Complete Purchase
                  </Button>
                </Grid>
              </Grid>
            </CardContent>
          </Card>

          <Card sx={{ mb: 4 }}>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Redeem Points
              </Typography>
              <Grid container spacing={2} alignItems="center">
                <Grid item xs={12} sm={8}>
                  <TextField
                    fullWidth
                    label="Points to Redeem"
                    type="number"
                    value={points}
                    onChange={(e) => setPoints(e.target.value)}
                    margin="normal"
                  />
                </Grid>
                <Grid item xs={12} sm={4}>
                  <Button
                    variant="contained"
                    color="secondary"
                    onClick={handleRedeemPoints}
                    fullWidth
                    disabled={!points || parseInt(points) <= 0 || parseInt(points) > balance}
                  >
                    Redeem Points
                  </Button>
                </Grid>
              </Grid>
            </CardContent>
          </Card>

          <Box sx={{ mt: 4, display: 'flex', gap: 2, flexDirection: 'column' }}>
            <Button
              variant="outlined"
              color="primary"
              onClick={() => navigate('/history')}
              fullWidth
            >
              View Transaction History
            </Button>
            
            <Button
              variant="outlined"
              color="error"
              onClick={handleLogout}
              fullWidth
            >
              Logout
            </Button>
          </Box>
        </Paper>
      </Box>
    </Container>
  );
};

export default Dashboard; 