function getApiBaseUrl() {
  const { protocol, hostname, port } = window.location;

  // http://localhost:8080 (the dev API)
  if (hostname === 'localhost' && port == 5173) {
    return 'http://localhost:8080';
  }

  // Otherwise, use same origin as the app
  return `${protocol}//${hostname}${port ? `:${port}` : ''}`;
}

export const API_BASE_URL = getApiBaseUrl()