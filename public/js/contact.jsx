function NavAuth() {
  const [loggedIn, setLoggedIn] = React.useState(!!localStorage.getItem('token'));
  const handleLogout = () => { localStorage.removeItem('token'); window.location.reload(); };
  return (
    <>
      {loggedIn ? (
        <><a href="/profile">Profile</a>{' '}<a href="#" onClick={handleLogout}>Logout</a></>
      ) : <a href="/login">Login</a>}
      {' '}<a href="/booking" className="btn btn-primary btn-sm">Book Now</a>
    </>
  );
}

function ContactForm() {
  const [form, setForm] = React.useState({ name: '', email: '', subject: '', message: '' });
  const [msg, setMsg] = React.useState('');

  const handleChange = (e) => setForm({...form, [e.target.name]: e.target.value});

  const handleSubmit = (e) => {
    e.preventDefault();
    setMsg('Thank you! Your message has been sent. We will get back to you soon.');
    setForm({ name: '', email: '', subject: '', message: '' });
  };

  return (
    <div>
      {msg && (
        <div style={{padding:'1rem',borderRadius:'var(--radius-lg)',marginBottom:'1rem',background:'#efe',color:'#060'}}>{msg}</div>
      )}
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Name</label>
          <input type="text" name="name" className="form-control" placeholder="Your name" value={form.name} onChange={handleChange} required />
        </div>
        <div className="form-group">
          <label>Email</label>
          <input type="email" name="email" className="form-control" placeholder="you@example.com" value={form.email} onChange={handleChange} required />
        </div>
        <div className="form-group">
          <label>Subject</label>
          <input type="text" name="subject" className="form-control" placeholder="Message subject" value={form.subject} onChange={handleChange} />
        </div>
        <div className="form-group">
          <label>Message</label>
          <textarea name="message" className="form-control" rows="5" placeholder="Your message..." value={form.message} onChange={handleChange} required></textarea>
        </div>
        <button type="submit" className="btn btn-primary btn-lg" style={{width:'100%'}}>Send</button>
      </form>
    </div>
  );
}

ReactDOM.createRoot(document.getElementById('navAuthContact')).render(<NavAuth />);
ReactDOM.createRoot(document.getElementById('contactFormRoot')).render(<ContactForm />);
