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

function RoomFilter() {
  const data = window.__ROOMS_DATA__ || {};
  const roomTypes = data.RoomTypes || [];
  const [active, setActive] = React.useState('all');

  React.useEffect(() => {
    document.querySelectorAll('.room-card').forEach(card => {
      if (active === 'all' || card.getAttribute('data-type') === String(active)) {
        card.style.display = '';
      } else {
        card.style.display = 'none';
      }
    });
  }, [active]);

  return (
    <div className="filter-bar">
      <button className={'filter-btn ' + (active==='all' ? 'active':'')} onClick={() => setActive('all')}>All Rooms</button>
      {roomTypes.map(rt => (
        <button key={rt.id} className={'filter-btn ' + (active===rt.id ? 'active':'')} onClick={() => setActive(rt.id)}>{rt.name}</button>
      ))}
    </div>
  );
}

ReactDOM.createRoot(document.getElementById('navAuthRooms')).render(<NavAuth />);
ReactDOM.createRoot(document.getElementById('roomFilterRoot')).render(<RoomFilter />);
