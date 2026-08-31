import { useEffect, useState } from 'react';
import {
  Box,
  Stack,
  Card,
  CardContent,
  CardActions,
  TextField,
  Button,
  Typography,
  IconButton,
  Alert,
  CircularProgress,
  Fade,
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import SaveIcon from '@mui/icons-material/Save';
import AddIcon from '@mui/icons-material/Add';
import { listNotes, createNote, updateNote, deleteNote } from '../api/client';

export default function NotesPage() {
  const [notes, setNotes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [newContent, setNewContent] = useState('');
  const [creating, setCreating] = useState(false);
  const [editedContent, setEditedContent] = useState({});
  const [savingId, setSavingId] = useState(null);
  const [deletingId, setDeletingId] = useState(null);

  const refresh = async () => {
    setLoading(true);
    setError('');
    try {
      const data = await listNotes();
      setNotes(data || []);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    refresh();
  }, []);

  const handleCreate = async () => {
    if (!newContent.trim()) return;
    setCreating(true);
    setError('');
    try {
      await createNote(newContent.trim());
      setNewContent('');
      await refresh();
    } catch (err) {
      setError(err.message);
    } finally {
      setCreating(false);
    }
  };

  const handleSave = async (id) => {
    const content = editedContent[id];
    if (content === undefined) return;
    setSavingId(id);
    setError('');
    try {
      await updateNote(id, content);
      await refresh();
    } catch (err) {
      setError(err.message);
    } finally {
      setSavingId(null);
    }
  };

  const handleDelete = async (id) => {
    setDeletingId(id);
    setError('');
    try {
      await deleteNote(id);
      setNotes((prev) => prev.filter((n) => n.id !== id));
    } catch (err) {
      setError(err.message);
    } finally {
      setDeletingId(null);
    }
  };

  return (
    <Box>
      <Typography variant="h5" fontWeight={600} gutterBottom>
        Notes
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
        These notes are the source material the AI quiz generator studies to write questions.
      </Typography>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError('')}>
          {error}
        </Alert>
      )}

      <Card variant="outlined" sx={{ mb: 3 }}>
        <CardContent>
          <TextField
            label="New note"
            placeholder="Write something you want to be quizzed on..."
            multiline
            minRows={2}
            fullWidth
            value={newContent}
            onChange={(e) => setNewContent(e.target.value)}
          />
        </CardContent>
        <CardActions sx={{ justifyContent: 'flex-end', px: 2, pb: 2 }}>
          <Button
            variant="contained"
            startIcon={creating ? <CircularProgress size={16} color="inherit" /> : <AddIcon />}
            onClick={handleCreate}
            disabled={creating || !newContent.trim()}
          >
            Add note
          </Button>
        </CardActions>
      </Card>

      {loading ? (
        <Box display="flex" justifyContent="center" py={4}>
          <CircularProgress />
        </Box>
      ) : notes.length === 0 ? (
        <Typography color="text.secondary">No notes yet. Add one above to get started.</Typography>
      ) : (
        <Stack spacing={2}>
          {notes.map((note) => (
            <Fade in key={note.id}>
              <Card variant="outlined">
                <CardContent>
                  <TextField
                    fullWidth
                    multiline
                    minRows={2}
                    value={editedContent[note.id] ?? note.content}
                    onChange={(e) =>
                      setEditedContent((prev) => ({ ...prev, [note.id]: e.target.value }))
                    }
                  />
                </CardContent>
                <CardActions sx={{ justifyContent: 'space-between', px: 2, pb: 2 }}>
                  <Typography variant="caption" color="text.secondary">
                    Updated {new Date(note.updatedAt).toLocaleString()}
                  </Typography>
                  <Box>
                    <Button
                      size="small"
                      startIcon={
                        savingId === note.id ? (
                          <CircularProgress size={14} color="inherit" />
                        ) : (
                          <SaveIcon />
                        )
                      }
                      onClick={() => handleSave(note.id)}
                      disabled={
                        savingId === note.id ||
                        editedContent[note.id] === undefined ||
                        editedContent[note.id] === note.content
                      }
                    >
                      Save
                    </Button>
                    <IconButton
                      color="error"
                      onClick={() => handleDelete(note.id)}
                      disabled={deletingId === note.id}
                    >
                      {deletingId === note.id ? <CircularProgress size={18} /> : <DeleteIcon />}
                    </IconButton>
                  </Box>
                </CardActions>
              </Card>
            </Fade>
          ))}
        </Stack>
      )}
    </Box>
  );
}
