function NavAuthProfile() {
  const handleLogout = () => { localStorage.removeItem('token'); window.location.href = '/'; };
  return <a href="#" className="btn btn-primary btn-sm" onClick={handleLogout}>Logout</a>;
}

function ProfileApp() {
  const [bookings, setBookings] = React.useState([]);
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState('');
  const [reviewForm, setReviewForm] = React.useState({ booking_id: '', rating: '5', comment: '' });
  const [reviewMsg, setReviewMsg] = React.useState('');
  const [reviewMsgType, setReviewMsgType] = React.useState('');

  const token = localStorage.getItem('token');

  React.useEffect(() => {
    if (!token) { window.location.href = '/login'; return; }
    fetch('/api/booking/all', { headers: { 'Authorization': 'Bearer ' + token } })
      .then(r => r.json())
      .then(data => {
        const list = Array.isArray(data) ? data : (data.bookings || []);
        setBookings(list);
        setLoading(false);
      })
      .catch(err => { setError(err.message); setLoading(false); });
  }, []);

  const cancelBooking = async (id) => {
    const res = await fetch('/api/booking/delete?id=' + id, {
      method: 'DELETE', headers: { 'Authorization': 'Bearer ' + token }
    });
    if (res.ok) setBookings(prev => prev.filter(b => (b.id || b.ID) !== id));
    else alert('Error canceling');
  };

  const handleReviewChange = (e) => setReviewForm({...reviewForm, [e.target.name]: e.target.value});

  const submitReview = async (e) => {
    e.preventDefault();
    setReviewMsg('');
    try {
      const res = await fetch('/api/reviews/create', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'Authorization': 'Bearer ' + token },
        body: JSON.stringify({
          booking_id: parseInt(reviewForm.booking_id),
          rating: parseInt(reviewForm.rating),
          comment: reviewForm.comment
        })
      });
      const data = await res.json();
      if (res.ok) { setReviewMsg('Review submitted!'); setReviewMsgType('success'); }
      else { setReviewMsg('Error: ' + (data.error || JSON.stringify(data))); setReviewMsgType('error'); }
    } catch(err) { setReviewMsg('Network error: ' + err.message); setReviewMsgType('error'); }
  };

  const thStyle = { padding: '.75rem', textAlign: 'left' };
  const tdStyle = { padding: '.75rem' };

  return (
    <div>
      <div className="card" style={{padding:'2rem'}}>
        <h2 style={{marginBottom:'1.5rem'}}>My Bookings</h2>
        {error && <div style={{padding:'1rem',borderRadius:'var(--radius-lg)',marginBottom:'1rem',background:'#fee',color:'#c00'}}>Error: {error}</div>}
        {loading ? <p>Loading...</p> : bookings.length === 0 ? (
          <p style={{color:'var(--color-text-secondary)',textAlign:'center',padding:'2rem'}}>You have no bookings yet.</p>
        ) : (
          <div style={{overflowX:'auto'}}>
            <table style={{width:'100%',borderCollapse:'collapse'}}>
              <thead>
                <tr style={{borderBottom:'2px solid var(--color-border)'}}>
                  <th style={thStyle}>ID</th><th style={thStyle}>Room</th>
                  <th style={thStyle}>Check-in</th><th style={thStyle}>Check-out</th>
                  <th style={thStyle}>Total</th><th style={thStyle}>Status</th>
                  <th style={thStyle}>Actions</th>
                </tr>
              </thead>
              <tbody>
                {bookings.map(b => {
                  const id = b.id || b.ID;
                  const ci = (b.check_in || b.start_date || '').substring(0,10);
                  const co = (b.check_out || b.end_date || '').substring(0,10);
                  return (
                    <tr key={id} style={{borderBottom:'1px solid var(--color-border)'}}>
                      <td style={tdStyle}>{id}</td>
                      <td style={tdStyle}>{b.room_id || b.RoomID}</td>
                      <td style={tdStyle}>{ci}</td>
                      <td style={tdStyle}>{co}</td>
                      <td style={tdStyle}>${b.total_price || b.TotalPrice || 0}</td>
                      <td style={tdStyle}>{b.status || ''}</td>
                      <td style={tdStyle}>
                        <button className="btn btn-outline btn-sm" onClick={() => cancelBooking(id)}>Cancel</button>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>

      <div className="card" style={{padding:'2rem', marginTop:'2rem'}}>
        <h2 style={{marginBottom:'1.5rem'}}>Leave a Review</h2>
        {reviewMsg && (
          <div style={{padding:'1rem',borderRadius:'var(--radius-lg)',marginBottom:'1rem',
            background: reviewMsgType==='success' ? '#efe':'#fee',
            color: reviewMsgType==='success' ? '#060':'#c00'}}>{reviewMsg}</div>
        )}
        <form onSubmit={submitReview}>
          <div className="form-group">
            <label>Booking ID</label>
            <input type="number" name="booking_id" className="form-control" placeholder="Enter ID" value={reviewForm.booking_id} onChange={handleReviewChange} required />
          </div>
          <div className="form-group">
            <label>Rating (1-5)</label>
            <select name="rating" className="form-control" value={reviewForm.rating} onChange={handleReviewChange} required>
              <option value="5">5 — Excellent</option>
              <option value="4">4 — Good</option>
              <option value="3">3 — Average</option>
              <option value="2">2 — Poor</option>
              <option value="1">1 — Terrible</option>
            </select>
          </div>
          <div className="form-group">
            <label>Comment</label>
            <textarea name="comment" className="form-control" rows="3" placeholder="Your review..." value={reviewForm.comment} onChange={handleReviewChange}></textarea>
          </div>
          <button type="submit" className="btn btn-primary">Submit Review</button>
        </form>
      </div>
    </div>
  );
}

ReactDOM.createRoot(document.getElementById('navAuthProfile')).render(<NavAuthProfile />);
ReactDOM.createRoot(document.getElementById('profileRoot')).render(<ProfileApp />);
