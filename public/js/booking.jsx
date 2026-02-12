function NavAuth() {
  const [loggedIn, setLoggedIn] = React.useState(!!localStorage.getItem('token'));
  const handleLogout = () => { localStorage.removeItem('token'); window.location.reload(); };
  if (loggedIn) {
    return <><a href="/profile">Profile</a>{' '}<a href="#" onClick={handleLogout}>Logout</a></>;
  }
  return <a href="/login">Login</a>;
}

// Decode JWT to get user_id
function parseJwt(token) {
  try {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(c => {
      return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));
    return JSON.parse(jsonPayload);
  } catch (e) {
    return null;
  }
}

function BookingForm() {
  const data = window.__BOOKING_DATA__ || {};
  const rooms = (data.Rooms || []).filter(r => r.status !== 'booked');
  const mealPlans = data.MealPlans || [];
  const packages = data.Packages || [];
  const services = data.Services || [];

  const [form, setForm] = React.useState({
    room_id: '', check_in: '', check_out: '',
    meal_plan_id: '0', package_id: '0', payment_method: 'card'
  });
  const [selSvc, setSelSvc] = React.useState([]);
  const [msg, setMsg] = React.useState('');
  const [msgType, setMsgType] = React.useState('');
  const [loading, setLoading] = React.useState(false);

  const handleChange = (e) => setForm({...form, [e.target.name]: e.target.value});
  const toggleSvc = (id) => setSelSvc(prev => prev.includes(id) ? prev.filter(s => s !== id) : [...prev, id]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const token = localStorage.getItem('token');
    if (!token) { setMsg('Please login to book.'); setMsgType('error'); return; }
    
    // Get user_id from token
    const payload = parseJwt(token);
    if (!payload || !payload.id) { setMsg('Authorization error. Please login again.'); setMsgType('error'); return; }
    const userId = payload.id;
    
    setLoading(true); setMsg('');
    try {
      const res = await fetch('/api/booking', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'Authorization': 'Bearer ' + token },
        body: JSON.stringify({
          user_id: userId, room_id: parseInt(form.room_id),
          start_date: form.check_in, end_date: form.check_out,
          mealplan_id: parseInt(form.meal_plan_id) || null,
          package_id: parseInt(form.package_id) || null,
          payment_method: form.payment_method,
          total_price: 0
        })
      });
      const result = await res.json();
      if (res.ok) { setMsg('Booking created! ID: ' + (result.id || '') + '. Payment: ' + form.payment_method); setMsgType('success'); }
      else { setMsg('Error: ' + (result.error || JSON.stringify(result))); setMsgType('error'); }
    } catch(err) { setMsg('Network error: ' + err.message); setMsgType('error'); }
    setLoading(false);
  };

  const msgStyle = { padding:'1rem', borderRadius:'var(--radius-lg)', marginBottom:'1rem' };

  return (
    <div>
      {msg && (
        <div style={{...msgStyle, background: msgType==='success' ? '#efe':'#fee', color: msgType==='success' ? '#060':'#c00'}}>
          {msg}
        </div>
      )}
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Room</label>
          <select name="room_id" className="form-control" value={form.room_id} onChange={handleChange} required>
            <option value="">— Select a room —</option>
            {rooms.map(r => (
              <option key={r.id} value={r.id}>#{r.code} — {r.type ? r.type.name : ''} — ${r.price}/night</option>
            ))}
          </select>
        </div>
        <div style={{display:'grid', gridTemplateColumns:'1fr 1fr', gap:'1rem'}}>
          <div className="form-group">
            <label>Check-in Date</label>
            <input type="date" name="check_in" className="form-control" value={form.check_in} onChange={handleChange} required />
          </div>
          <div className="form-group">
            <label>Check-out Date</label>
            <input type="date" name="check_out" className="form-control" value={form.check_out} onChange={handleChange} required />
          </div>
        </div>
        <div className="form-group">
          <label>Meal Plan</label>
          <select name="meal_plan_id" className="form-control" value={form.meal_plan_id} onChange={handleChange}>
            <option value="0">No Meals</option>
            {mealPlans.map(m => (
              <option key={m.id} value={m.id}>{m.name} — ${m.price_per_day}/day</option>
            ))}
          </select>
        </div>
        <div className="form-group">
          <label>Package</label>
          <select name="package_id" className="form-control" value={form.package_id} onChange={handleChange}>
            <option value="0">No Package</option>
            {packages.map(p => (
              <option key={p.id} value={p.id}>{p.name} ({p.price_modifier * 100}% discount)</option>
            ))}
          </select>
        </div>
        <div className="form-group">
          <label>Additional Services</label>
          {services.map(s => (
            <div key={s.id} style={{margin:'.5rem 0'}}>
              <label style={{display:'flex', alignItems:'center', gap:'.5rem', cursor:'pointer'}}>
                <input type="checkbox" checked={selSvc.includes(s.id)} onChange={() => toggleSvc(s.id)} />
                {s.name} — ${s.price}
              </label>
            </div>
          ))}
        </div>
        <div className="form-group">
          <label>Payment Method</label>
          <select name="payment_method" className="form-control" value={form.payment_method} onChange={handleChange} required>
            <option value="card">Credit Card</option>
            <option value="cash">Cash on Check-in</option>
            <option value="online">Online Transfer</option>
          </select>
        </div>
        <button type="submit" className="btn btn-primary btn-lg" style={{width:'100%', marginTop:'1rem'}} disabled={loading}>
          {loading ? 'Submitting...' : 'Book Now'}
        </button>
      </form>
    </div>
  );
}

const navEl = document.getElementById('navAuth');
if (navEl) ReactDOM.createRoot(navEl).render(<NavAuth />);
ReactDOM.createRoot(document.getElementById('bookingRoot')).render(<BookingForm />);
