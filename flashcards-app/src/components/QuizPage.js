import { useEffect, useRef, useState } from 'react';
import {
  Box,
  Stack,
  Paper,
  TextField,
  Button,
  Typography,
  Avatar,
  CircularProgress,
  Alert,
  IconButton,
  Tooltip,
} from '@mui/material';
import SendIcon from '@mui/icons-material/Send';
import RestartAltIcon from '@mui/icons-material/RestartAlt';
import SchoolIcon from '@mui/icons-material/School';
import PersonIcon from '@mui/icons-material/Person';
import { generateQuiz } from '../api/client';

export default function QuizPage() {
  const [conversationId, setConversationId] = useState('');
  const [messages, setMessages] = useState([]);
  const [answer, setAnswer] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [started, setStarted] = useState(false);
  const bottomRef = useRef(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const start = async () => {
    setLoading(true);
    setError('');
    try {
      const res = await generateQuiz('', []);
      setConversationId(res.conversationId);
      setMessages(res.messages || []);
      setStarted(true);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const sendAnswer = async () => {
    if (!answer.trim()) return;
    const nextMessages = [...messages, { role: 'user', content: answer.trim() }];
    setMessages(nextMessages);
    setAnswer('');
    setLoading(true);
    setError('');
    try {
      const res = await generateQuiz(conversationId, nextMessages);
      setConversationId(res.conversationId);
      setMessages(res.messages || []);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const reset = () => {
    setConversationId('');
    setMessages([]);
    setAnswer('');
    setError('');
    setStarted(false);
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      sendAnswer();
    }
  };

  return (
    <Box>
      <Stack direction="row" justifyContent="space-between" alignItems="center" sx={{ mb: 2 }}>
        <Box>
          <Typography variant="h5" fontWeight={600}>
            Quiz
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Answer AI-generated questions based on your notes.
          </Typography>
        </Box>
        {started && (
          <Tooltip title="Restart quiz">
            <IconButton onClick={reset}>
              <RestartAltIcon />
            </IconButton>
          </Tooltip>
        )}
      </Stack>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError('')}>
          {error}
        </Alert>
      )}

      {!started ? (
        <Box display="flex" flexDirection="column" alignItems="center" py={6}>
          <SchoolIcon sx={{ fontSize: 56, color: 'primary.main', mb: 2 }} />
          <Typography sx={{ mb: 3 }} color="text.secondary">
            Ready to test what you've learned from your notes?
          </Typography>
          <Button
            variant="contained"
            size="large"
            onClick={start}
            disabled={loading}
            startIcon={loading ? <CircularProgress size={18} color="inherit" /> : null}
          >
            Start quiz
          </Button>
        </Box>
      ) : (
        <>
          <Paper
            variant="outlined"
            sx={{
              p: 2,
              mb: 2,
              maxHeight: 480,
              overflowY: 'auto',
              bgcolor: 'grey.50',
            }}
          >
            <Stack spacing={2}>
              {messages.map((m, i) => (
                <Stack
                  key={i}
                  direction="row"
                  spacing={1.5}
                  justifyContent={m.role === 'user' ? 'flex-end' : 'flex-start'}
                >
                  {m.role !== 'user' && (
                    <Avatar sx={{ bgcolor: 'primary.main', width: 32, height: 32 }}>
                      <SchoolIcon fontSize="small" />
                    </Avatar>
                  )}
                  <Paper
                    elevation={0}
                    sx={{
                      p: 1.5,
                      maxWidth: '75%',
                      bgcolor: m.role === 'user' ? 'primary.main' : 'background.paper',
                      color: m.role === 'user' ? 'primary.contrastText' : 'text.primary',
                      borderRadius: 2,
                    }}
                  >
                    <Typography variant="body2" whiteSpace="pre-wrap">
                      {m.content}
                    </Typography>
                  </Paper>
                  {m.role === 'user' && (
                    <Avatar sx={{ bgcolor: 'secondary.main', width: 32, height: 32 }}>
                      <PersonIcon fontSize="small" />
                    </Avatar>
                  )}
                </Stack>
              ))}
              {loading && (
                <Stack direction="row" spacing={1.5} alignItems="center">
                  <Avatar sx={{ bgcolor: 'primary.main', width: 32, height: 32 }}>
                    <SchoolIcon fontSize="small" />
                  </Avatar>
                  <CircularProgress size={20} />
                </Stack>
              )}
              <div ref={bottomRef} />
            </Stack>
          </Paper>

          <Stack direction="row" spacing={1}>
            <TextField
              fullWidth
              placeholder="Type your answer..."
              value={answer}
              onChange={(e) => setAnswer(e.target.value)}
              onKeyDown={handleKeyDown}
              disabled={loading}
              multiline
              maxRows={4}
            />
            <Button
              variant="contained"
              endIcon={<SendIcon />}
              onClick={sendAnswer}
              disabled={loading || !answer.trim()}
            >
              Send
            </Button>
          </Stack>
        </>
      )}
    </Box>
  );
}
