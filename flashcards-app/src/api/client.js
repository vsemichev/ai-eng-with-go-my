const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || '';

async function request(path, options = {}) {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    headers: { 'Content-Type': 'application/json', ...options.headers },
    ...options,
  });

  if (!response.ok) {
    let message = `Request failed with status ${response.status}`;
    try {
      const body = await response.json();
      if (body && body.error) message = body.error;
    } catch {
      // response had no JSON body; keep default message
    }
    throw new Error(message);
  }

  if (response.status === 204) return null;
  return response.json();
}

// ----- Notes -----
export function listNotes() {
  return request('/notes');
}

export function getNote(id) {
  return request(`/notes/${id}`);
}

export function createNote(content) {
  return request('/notes', {
    method: 'POST',
    body: JSON.stringify({ content }),
  });
}

export function updateNote(id, content) {
  return request(`/notes/${id}`, {
    method: 'PUT',
    body: JSON.stringify({ content }),
  });
}

export function deleteNote(id) {
  return request(`/notes/${id}`, { method: 'DELETE' });
}

// ----- Quiz -----
export function generateQuiz(conversationId, messages) {
  return request('/quiz/generate', {
    method: 'POST',
    body: JSON.stringify({ conversationId: conversationId || '', messages }),
  });
}

// ----- Health -----
export function checkHealth() {
  return request('/health');
}
