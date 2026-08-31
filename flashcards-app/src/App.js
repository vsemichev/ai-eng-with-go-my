import { useEffect, useState } from 'react';
import {
  ThemeProvider,
  createTheme,
  CssBaseline,
  AppBar,
  Toolbar,
  Typography,
  Container,
  Tabs,
  Tab,
  Box,
  Chip,
} from '@mui/material';
import NoteAltIcon from '@mui/icons-material/NoteAlt';
import QuizIcon from '@mui/icons-material/Quiz';
import NotesPage from './components/NotesPage';
import QuizPage from './components/QuizPage';
import { checkHealth } from './api/client';

const theme = createTheme({
  palette: {
    mode: 'light',
    primary: { main: '#4f46e5' },
    secondary: { main: '#0ea5e9' },
    background: { default: '#f5f6fa' },
  },
  shape: { borderRadius: 12 },
});

function App() {
  const [tab, setTab] = useState(0);
  const [apiStatus, setApiStatus] = useState('checking');

  useEffect(() => {
    checkHealth()
      .then(() => setApiStatus('online'))
      .catch(() => setApiStatus('offline'));
  }, []);

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <AppBar position="static" color="transparent" elevation={0} sx={{ bgcolor: 'background.paper' }}>
        <Toolbar>
          <Typography variant="h6" fontWeight={700} sx={{ flexGrow: 1 }}>
            📚 Flashcards
          </Typography>
          <Chip
            size="small"
            label={apiStatus === 'online' ? 'API online' : apiStatus === 'offline' ? 'API offline' : 'Checking...'}
            color={apiStatus === 'online' ? 'success' : apiStatus === 'offline' ? 'error' : 'default'}
          />
        </Toolbar>
        <Container maxWidth="md">
          <Tabs value={tab} onChange={(_, v) => setTab(v)}>
            <Tab icon={<NoteAltIcon />} iconPosition="start" label="Notes" />
            <Tab icon={<QuizIcon />} iconPosition="start" label="Quiz" />
          </Tabs>
        </Container>
      </AppBar>

      <Container maxWidth="md" sx={{ py: 4 }}>
        <Box role="tabpanel" hidden={tab !== 0}>
          {tab === 0 && <NotesPage />}
        </Box>
        <Box role="tabpanel" hidden={tab !== 1}>
          {tab === 1 && <QuizPage />}
        </Box>
      </Container>
    </ThemeProvider>
  );
}

export default App;
