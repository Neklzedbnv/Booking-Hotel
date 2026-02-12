function RegisterForm() {
  const [form, setForm] = React.useState({
    name: '', email: '', phone: '', password: '', password2: ''
  });
  const [msg, setMsg] = React.useState('');
  const [msgType, setMsgType] = React.useState('');
  const [loading, setLoading] = React.useState(false);

  const handleChange = (e) => {
    setForm({...form, [e.target.name]: e.target.value});
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setMsg('');
    if (form.password !== form.password2) {
      setMsg('Passwords don\'t match!');
      setMsgType('error');
      return;
    }
    setLoading(true);
    try {
      const res = await fetch('/api/auth/register', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
          fullname: form.name, email: form.email,
          phone: form.phone, password: form.password
        })
      });
      const data = await res.json();
      if (res.ok) {
        setMsg('Registration successful! Redirecting to login...');
        setMsgType('success');
        setTimeout(() => { window.location.href = '/login'; }, 1000);
      } else {
        setMsg('Error: ' + (data.error || JSON.stringify(data)));
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
          <label>Name</label>
          <input type="text" name="name" className="form-control" placeholder="Your name"
            value={form.name} onChange={handleChange} required />
        </div>
        <div className="form-group">
          <label>Email</label>
          <input type="email" name="email" className="form-control" placeholder="you@example.com"
            value={form.email} onChange={handleChange} required />
        </div>
        <div className="form-group">
          <label>Phone</label>
          <input type="tel" name="phone" className="form-control" placeholder="+7 (777) 123-4567"
            value={form.phone} onChange={handleChange} />
        </div>
        <div className="form-group">
          <label>Password</label>
          <input type="password" name="password" className="form-control" placeholder="Minimum 6 characters"
            value={form.password} onChange={handleChange} required minLength="6" />
        </div>
        <div className="form-group">
          <label>Confirm Password</label>
          <input type="password" name="password2" className="form-control" placeholder="Repeat password"
            value={form.password2} onChange={handleChange} required />
        </div>
        <button type="submit" className="btn btn-primary btn-lg"
          style={{width: '100%', marginTop: '1rem'}} disabled={loading}>
          {loading ? 'Registering...' : 'Sign up'}
        </button>
      </form>
    </div>
  );
}
ReactDOM.createRoot(document.getElementById('regRoot')).render(<RegisterForm />);
