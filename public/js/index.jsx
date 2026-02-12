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
ReactDOM.createRoot(document.getElementById('navAuthIndex')).render(<NavAuth />);
