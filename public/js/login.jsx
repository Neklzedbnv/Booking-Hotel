const ADMIN_EMAIL = 'abzalbahktiarow2006@gmail.com';

function LoginForm() {
  const [email, setEmail] = React.useState('');
  const [password, setPassword] = React.useState('');
  const [msg, setMsg] = React.useState('');
  const [msgType, setMsgType] = React.useState('');
  const [loading, setLoading] = React.useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setMsg('');
    try {
      const res = await fetch('/api/auth/login', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({email, password})
      });
      const data = await res.json();
      if (res.ok && data.token) {
        localStorage.setItem('token', data.token);
        setMsg('Login successful! Redirecting...');
        setMsgType('success');
        // If admin - redirect to admin panel
        const redirectUrl = email === ADMIN_EMAIL ? '/admin' : '/profile';
        setTimeout(() => { window.location.href = redirectUrl; }, 800);
      } else {
        setMsg('Error: ' + (data.error || 'Invalid credentials'));
        setMsgType('error');
      }
    } catch(err) {
      setMsg('Network error: ' + err.message);
      setMsgType('error');
    }
    setLoading(false);
  };

  return (
    <div>
      {msg && (
        <div style={{
          padding: '1rem', borderRadius: 'var(--radius-lg)', marginBottom: '1rem',
          background: msgType === 'success' ? '#efe' : '#fee',
          color: msgType === 'success' ? '#060' : '#c00'
        }}>{msg}</div>
      )}
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Email</label>
          <input type="email" className="form-control" placeholder="you@example.com"
            value={email} onChange={e => setEmail(e.target.value)} required />
        </div>
        <div className="form-group">
          <label>Password</label>
          <input type="password" className="form-control" placeholder="••••••••"
            value={password} onChange={e => setPassword(e.target.value)} required />
        </div>
        <button type="submit" className="btn btn-primary btn-lg"
          style={{width: '100%', marginTop: '1rem'}} disabled={loading}>
          {loading ? 'Logging in...' : 'Login'}
        </button>
      </form>
    </div>
  );
}
ReactDOM.createRoot(document.getElementById('loginRoot')).render(<LoginForm />);
