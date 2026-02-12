const { useState, useEffect } = React;

// Decode JWT token
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

// Check admin
const ADMIN_EMAIL = 'abzalbahktiarow2006@gmail.com';

function checkAdmin() {
  const token = localStorage.getItem('token');
  if (!token) return false;
  const payload = parseJwt(token);
  if (!payload || !payload.email) return false;
  return payload.email === ADMIN_EMAIL;
}

// API helper
const api = {
  get: async (url) => {
    const token = localStorage.getItem('token');
    const res = await fetch(url, {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    return res.json();
  },
  post: async (url, data) => {
    const token = localStorage.getItem('token');
    const res = await fetch(url, {
      method: 'POST',
      headers: { 
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}` 
      },
      body: JSON.stringify(data)
    });
    return res.json();
  },
  delete: async (url) => {
    const token = localStorage.getItem('token');
    const res = await fetch(url, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${token}` }
    });
    return res.json();
  }
};

// Dashboard Section
function Dashboard({ stats }) {
  return (
    <div>
      <div className="stats-grid">
        <div className="stat-card">
          <h3>Users</h3>
          <div className="value">{stats.users_count || 0}</div>
        </div>
        <div className="stat-card">
          <h3>Rooms</h3>
          <div className="value">{stats.rooms_count || 0}</div>
        </div>
        <div className="stat-card">
          <h3>Bookings</h3>
          <div className="value">{stats.bookings_count || 0}</div>
        </div>
        <div className="stat-card">
          <h3>Total Revenue</h3>
          <div className="value">${(stats.total_revenue || 0).toFixed(2)}</div>
        </div>
      </div>
      
      {stats.bookings_by_status && (
        <div className="admin-section active">
          <h2>Booking Statistics</h2>
          <div className="stats-grid">
            {Object.entries(stats.bookings_by_status).map(([status, count]) => (
              <div key={status} className="stat-card">
                <h3>{statusLabels[status] || status}</h3>
                <div className="value">{count}</div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

const statusLabels = {
  'pending': 'Pending',
  'confirmed': 'Confirmed',
  'cancelled': 'Cancelled',
  'completed': 'Completed'
};

// Rooms Section
function RoomsSection({ rooms, roomTypes, onRefresh }) {
  const [showModal, setShowModal] = useState(false);
  const [editRoom, setEditRoom] = useState(null);
  const [form, setForm] = useState({ code: '', type_id: '', capacity: '', price: '' });

  const handleSubmit = async (e) => {
    e.preventDefault();
    const data = {
      code: form.code,
      type_id: parseInt(form.type_id),
      capacity: parseInt(form.capacity),
      price: parseFloat(form.price)
    };
    
    if (editRoom) {
      await api.post('/api/rooms/update?id=' + editRoom.id, data);
    } else {
      await api.post('/api/rooms/create', data);
    }
    
    setShowModal(false);
    setEditRoom(null);
    setForm({ code: '', type_id: '', capacity: '', price: '' });
    onRefresh();
  };

  const handleEdit = (room) => {
    setEditRoom(room);
    setForm({
      code: room.code,
      type_id: room.type_id.toString(),
      capacity: room.capacity.toString(),
      price: room.price.toString()
    });
    setShowModal(true);
  };

  const handleDelete = async (id) => {
    if (confirm('Delete room?')) {
      await api.delete('/api/rooms/delete?id=' + id);
      onRefresh();
    }
  };

  const handleStatusChange = async (id, status) => {
    await api.post('/api/rooms/update?id=' + id, { status });
    onRefresh();
  };

  return (
    <div className="admin-section active">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
        <h2>Room Management</h2>
        <button className="btn btn-primary btn-sm" onClick={() => setShowModal(true)}>+ Add Room</button>
      </div>
      
      <table className="admin-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Code</th>
            <th>Type</th>
            <th>Capacity</th>
            <th>Price</th>
            <th>Status</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {rooms.map(room => (
            <tr key={room.id}>
              <td>{room.id}</td>
              <td>{room.code}</td>
              <td>{roomTypes.find(t => t.id === room.type_id)?.name || room.type_id}</td>
              <td>{room.capacity}</td>
              <td>${room.price}</td>
              <td>
                <select 
                  value={room.status} 
                  onChange={(e) => handleStatusChange(room.id, e.target.value)}
                  className={`badge ${room.status === 'available' ? 'badge-success' : 'badge-warning'}`}
                  style={{ border: 'none', cursor: 'pointer' }}
                >
                  <option value="available">Available</option>
                  <option value="booked">Booked</option>
                  <option value="maintenance">Maintenance</option>
                  <option value="cleaning">Cleaning</option>
                </select>
              </td>
              <td>
                <button className="btn btn-outline btn-sm" onClick={() => handleEdit(room)}>✏️</button>
                <button className="btn btn-outline btn-sm" onClick={() => handleDelete(room.id)} style={{ marginLeft: '0.5rem', color: '#ef4444' }}>🗑️</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {showModal && (
        <div className="modal" onClick={() => setShowModal(false)}>
          <div className="modal-content" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3>{editRoom ? 'Edit Room' : 'Add Room'}</h3>
              <button className="modal-close" onClick={() => setShowModal(false)}>&times;</button>
            </div>
            <form className="admin-form" onSubmit={handleSubmit}>
              <div className="form-group">
                <label>Room Code</label>
                <input type="text" value={form.code} onChange={e => setForm({...form, code: e.target.value})} required />
              </div>
              <div className="form-group">
                <label>Room Type</label>
                <select value={form.type_id} onChange={e => setForm({...form, type_id: e.target.value})} required>
                  <option value="">Select type</option>
                  {roomTypes.map(t => <option key={t.id} value={t.id}>{t.name}</option>)}
                </select>
              </div>
              <div className="form-group">
                <label>Capacity</label>
                <input type="number" value={form.capacity} onChange={e => setForm({...form, capacity: e.target.value})} required />
              </div>
              <div className="form-group">
                <label>Price per Night ($)</label>
                <input type="number" step="0.01" value={form.price} onChange={e => setForm({...form, price: e.target.value})} required />
              </div>
              <button type="submit" className="btn btn-primary">Save</button>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

// Room Types Section
function RoomTypesSection({ roomTypes, onRefresh }) {
  const [showModal, setShowModal] = useState(false);
  const [editType, setEditType] = useState(null);
  const [form, setForm] = useState({ name: '', capacity: '', base_price: '' });

  const handleSubmit = async (e) => {
    e.preventDefault();
    const data = {
      name: form.name,
      capacity: parseInt(form.capacity),
      base_price: parseFloat(form.base_price)
    };
    
    if (editType) {
      await api.post('/api/admin/room-types/update', { ...data, id: editType.id });
    } else {
      await api.post('/api/rooms/types/create', data);
    }
    
    setShowModal(false);
    setEditType(null);
    setForm({ name: '', capacity: '', base_price: '' });
    onRefresh();
  };

  const handleEdit = (type) => {
    setEditType(type);
    setForm({
      name: type.name,
      capacity: type.capacity.toString(),
      base_price: type.base_price.toString()
    });
    setShowModal(true);
  };

  const handleDelete = async (id) => {
    if (confirm('Delete room type? All rooms of this type will also be deleted!')) {
      await api.delete('/api/admin/room-types/delete?id=' + id);
      onRefresh();
    }
  };

  return (
    <div className="admin-section active">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
        <h2>Room Types</h2>
        <button className="btn btn-primary btn-sm" onClick={() => setShowModal(true)}>+ Add Type</button>
      </div>
      
      <table className="admin-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Capacity</th>
            <th>Base Price</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {roomTypes.map(type => (
            <tr key={type.id}>
              <td>{type.id}</td>
              <td>{type.name}</td>
              <td>{type.capacity}</td>
              <td>${type.base_price}</td>
              <td>
                <button className="btn btn-outline btn-sm" onClick={() => handleEdit(type)}>✏️</button>
                <button className="btn btn-outline btn-sm" onClick={() => handleDelete(type.id)} style={{ marginLeft: '0.5rem', color: '#ef4444' }}>🗑️</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {showModal && (
        <div className="modal" onClick={() => setShowModal(false)}>
          <div className="modal-content" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3>{editType ? 'Edit Type' : 'Add Type'}</h3>
              <button className="modal-close" onClick={() => setShowModal(false)}>&times;</button>
            </div>
            <form className="admin-form" onSubmit={handleSubmit}>
              <div className="form-group">
                <label>Name</label>
                <input type="text" value={form.name} onChange={e => setForm({...form, name: e.target.value})} required />
              </div>
              <div className="form-group">
                <label>Capacity</label>
                <input type="number" value={form.capacity} onChange={e => setForm({...form, capacity: e.target.value})} required />
              </div>
              <div className="form-group">
                <label>Base Price ($)</label>
                <input type="number" step="0.01" value={form.base_price} onChange={e => setForm({...form, base_price: e.target.value})} required />
              </div>
              <button type="submit" className="btn btn-primary">Save</button>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

// Bookings Section
function BookingsSection({ bookings, onRefresh }) {
  const handleStatusChange = async (id, status) => {
    await api.post('/api/admin/bookings/status', { booking_id: id, status });
    onRefresh();
  };

  const statusBadge = (status) => {
    const classes = {
      'pending': 'badge-warning',
      'confirmed': 'badge-success',
      'cancelled': 'badge-danger',
      'completed': 'badge-info'
    };
    return classes[status] || 'badge-info';
  };

  return (
    <div className="admin-section active">
      <h2>Bookings</h2>
      <table className="admin-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Guest</th>
            <th>Email</th>
            <th>Room</th>
            <th>Dates</th>
            <th>Total</th>
            <th>Status</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {bookings.map(b => (
            <tr key={b.id}>
              <td>{b.id}</td>
              <td>{b.user_name}</td>
              <td>{b.user_email}</td>
              <td>{b.room_code}</td>
              <td>{new Date(b.start_date).toLocaleDateString()} - {new Date(b.end_date).toLocaleDateString()}</td>
              <td>${b.total_price}</td>
              <td><span className={`badge ${statusBadge(b.status)}`}>{statusLabels[b.status] || b.status}</span></td>
              <td>
                <select 
                  value={b.status} 
                  onChange={(e) => handleStatusChange(b.id, e.target.value)}
                  style={{ padding: '0.25rem', borderRadius: '4px', border: '1px solid #ddd' }}
                >
                  <option value="pending">Pending</option>
                  <option value="confirmed">Confirm</option>
                  <option value="cancelled">Cancel</option>
                  <option value="completed">Complete</option>
                </select>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

// Calendar Section
function CalendarSection({ bookings }) {
  const [currentMonth, setCurrentMonth] = useState(new Date());
  
  const getDaysInMonth = (date) => {
    const year = date.getFullYear();
    const month = date.getMonth();
    const firstDay = new Date(year, month, 1);
    const lastDay = new Date(year, month + 1, 0);
    const days = [];
    
    // Add empty days for alignment
    for (let i = 0; i < firstDay.getDay(); i++) {
      days.push(null);
    }
    
    for (let i = 1; i <= lastDay.getDate(); i++) {
      days.push(new Date(year, month, i));
    }
    
    return days;
  };

  const getBookingsForDay = (date) => {
    if (!date) return [];
    return bookings.filter(b => {
      const start = new Date(b.start_date);
      const end = new Date(b.end_date);
      return date >= start && date <= end;
    });
  };

  const days = getDaysInMonth(currentMonth);
  const monthNames = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'];

  return (
    <div className="admin-section active">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
        <h2>Occupancy Calendar</h2>
        <div>
          <button className="btn btn-outline btn-sm" onClick={() => setCurrentMonth(new Date(currentMonth.setMonth(currentMonth.getMonth() - 1)))}>←</button>
          <span style={{ margin: '0 1rem' }}>{monthNames[currentMonth.getMonth()]} {currentMonth.getFullYear()}</span>
          <button className="btn btn-outline btn-sm" onClick={() => setCurrentMonth(new Date(currentMonth.setMonth(currentMonth.getMonth() + 1)))}>→</button>
        </div>
      </div>
      
      <div className="calendar-grid">
        {['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'].map(d => (
          <div key={d} className="calendar-day header">{d}</div>
        ))}
        {days.map((day, i) => (
          <div key={i} className="calendar-day">
            {day && (
              <>
                <div style={{ fontWeight: '600' }}>{day.getDate()}</div>
                {getBookingsForDay(day).map(b => (
                  <div key={b.id} className="calendar-booking" title={`${b.user_name} - ${b.room_code}`}>
                    {b.room_code}
                  </div>
                ))}
              </>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}

// Services Section
function ServicesSection({ services, onRefresh }) {
  const [showModal, setShowModal] = useState(false);
  const [editService, setEditService] = useState(null);
  const [form, setForm] = useState({ name: '', price: '' });

  const handleSubmit = async (e) => {
    e.preventDefault();
    const data = { name: form.name, price: parseFloat(form.price) };
    
    if (editService) {
      await api.post('/api/services/update?id=' + editService.id, data);
    } else {
      await api.post('/api/services/create', data);
    }
    
    setShowModal(false);
    setEditService(null);
    setForm({ name: '', price: '' });
    onRefresh();
  };

  const handleEdit = (svc) => {
    setEditService(svc);
    setForm({ name: svc.name, price: svc.price.toString() });
    setShowModal(true);
  };

  const handleDelete = async (id) => {
    if (confirm('Delete service?')) {
      await api.delete('/api/services/delete?id=' + id);
      onRefresh();
    }
  };

  return (
    <div className="admin-section active">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
        <h2>Services</h2>
        <button className="btn btn-primary btn-sm" onClick={() => setShowModal(true)}>+ Add Service</button>
      </div>
      
      <table className="admin-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Price</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {services.map(svc => (
            <tr key={svc.id}>
              <td>{svc.id}</td>
              <td>{svc.name}</td>
              <td>${svc.price}</td>
              <td>
                <button className="btn btn-outline btn-sm" onClick={() => handleEdit(svc)}>✏️</button>
                <button className="btn btn-outline btn-sm" onClick={() => handleDelete(svc.id)} style={{ marginLeft: '0.5rem', color: '#ef4444' }}>🗑️</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {showModal && (
        <div className="modal" onClick={() => setShowModal(false)}>
          <div className="modal-content" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3>{editService ? 'Edit Service' : 'Add Service'}</h3>
              <button className="modal-close" onClick={() => setShowModal(false)}>&times;</button>
            </div>
            <form className="admin-form" onSubmit={handleSubmit}>
              <div className="form-group">
                <label>Name</label>
                <input type="text" value={form.name} onChange={e => setForm({...form, name: e.target.value})} required />
              </div>
              <div className="form-group">
                <label>Price ($)</label>
                <input type="number" step="0.01" value={form.price} onChange={e => setForm({...form, price: e.target.value})} required />
              </div>
              <button type="submit" className="btn btn-primary">Save</button>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

// Packages Section
function PackagesSection({ packages, onRefresh }) {
  const [showModal, setShowModal] = useState(false);
  const [editPkg, setEditPkg] = useState(null);
  const [form, setForm] = useState({ name: '', description: '', price_modifier: '' });

  const handleSubmit = async (e) => {
    e.preventDefault();
    const data = { 
      name: form.name, 
      description: form.description,
      price_modifier: parseFloat(form.price_modifier) 
    };
    
    if (editPkg) {
      await api.post('/api/packages/update?id=' + editPkg.id, data);
    } else {
      await api.post('/api/packages/create', data);
    }
    
    setShowModal(false);
    setEditPkg(null);
    setForm({ name: '', description: '', price_modifier: '' });
    onRefresh();
  };

  const handleEdit = (pkg) => {
    setEditPkg(pkg);
    setForm({ 
      name: pkg.name, 
      description: pkg.description || '',
      price_modifier: pkg.price_modifier.toString() 
    });
    setShowModal(true);
  };

  const handleDelete = async (id) => {
    if (confirm('Delete package?')) {
      await api.delete('/api/packages/delete?id=' + id);
      onRefresh();
    }
  };

  return (
    <div className="admin-section active">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
        <h2>Packages</h2>
        <button className="btn btn-primary btn-sm" onClick={() => setShowModal(true)}>+ Add Package</button>
      </div>
      
      <table className="admin-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Description</th>
            <th>Price Modifier</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {packages.map(pkg => (
            <tr key={pkg.id}>
              <td>{pkg.id}</td>
              <td>{pkg.name}</td>
              <td>{pkg.description}</td>
              <td>${pkg.price_modifier}</td>
              <td>
                <button className="btn btn-outline btn-sm" onClick={() => handleEdit(pkg)}>✏️</button>
                <button className="btn btn-outline btn-sm" onClick={() => handleDelete(pkg.id)} style={{ marginLeft: '0.5rem', color: '#ef4444' }}>🗑️</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {showModal && (
        <div className="modal" onClick={() => setShowModal(false)}>
          <div className="modal-content" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3>{editPkg ? 'Edit Package' : 'Add Package'}</h3>
              <button className="modal-close" onClick={() => setShowModal(false)}>&times;</button>
            </div>
            <form className="admin-form" onSubmit={handleSubmit}>
              <div className="form-group">
                <label>Name</label>
                <input type="text" value={form.name} onChange={e => setForm({...form, name: e.target.value})} required />
              </div>
              <div className="form-group">
                <label>Description</label>
                <input type="text" value={form.description} onChange={e => setForm({...form, description: e.target.value})} />
              </div>
              <div className="form-group">
                <label>Price Modifier ($)</label>
                <input type="number" step="0.01" value={form.price_modifier} onChange={e => setForm({...form, price_modifier: e.target.value})} required />
              </div>
              <button type="submit" className="btn btn-primary">Save</button>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

// Users Section
function UsersSection({ users, onRefresh }) {
  const handleRoleChange = async (userId, role) => {
    await api.post('/api/admin/users/role', { user_id: userId, role });
    onRefresh();
  };

  const handleBlockToggle = async (userId, blocked) => {
    await api.post('/api/admin/users/block', { user_id: userId, blocked });
    onRefresh();
  };

  return (
    <div className="admin-section active">
      <h2>Users</h2>
      <table className="admin-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Email</th>
            <th>Role</th>
            <th>Status</th>
            <th>Registration Date</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {users.map(user => (
            <tr key={user.id}>
              <td>{user.id}</td>
              <td>{user.fullname}</td>
              <td>{user.email}</td>
              <td>
                <select 
                  value={user.role} 
                  onChange={(e) => handleRoleChange(user.id, e.target.value)}
                  style={{ padding: '0.25rem', borderRadius: '4px', border: '1px solid #ddd' }}
                >
                  <option value="user">User</option>
                  <option value="admin">Administrator</option>
                </select>
              </td>
              <td>
                <span className={`badge ${user.is_blocked ? 'badge-danger' : 'badge-success'}`}>
                  {user.is_blocked ? 'Blocked' : 'Active'}
                </span>
              </td>
              <td>{new Date(user.created_at).toLocaleDateString()}</td>
              <td>
                <button 
                  className={`btn btn-sm ${user.is_blocked ? 'btn-primary' : 'btn-outline'}`}
                  onClick={() => handleBlockToggle(user.id, !user.is_blocked)}
                  style={{ color: user.is_blocked ? '' : '#ef4444' }}
                >
                  {user.is_blocked ? '✓ Unblock' : '🚫 Block'}
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

// Reviews Section
function ReviewsSection({ reviews, onRefresh }) {
  const handleDelete = async (id) => {
    if (confirm('Delete review?')) {
      await api.delete('/api/reviews/delete?id=' + id);
      onRefresh();
    }
  };

  return (
    <div className="admin-section active">
      <h2>Reviews</h2>
      <table className="admin-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Booking</th>
            <th>Rating</th>
            <th>Comment</th>
            <th>Date</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {reviews.map(rev => (
            <tr key={rev.id}>
              <td>{rev.id}</td>
              <td>#{rev.booking_id}</td>
              <td>{'⭐'.repeat(rev.rating)}</td>
              <td>{rev.comment}</td>
              <td>{new Date(rev.created_at).toLocaleDateString()}</td>
              <td>
                <button className="btn btn-outline btn-sm" onClick={() => handleDelete(rev.id)} style={{ color: '#ef4444' }}>🗑️ Delete</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

// Images Section
function ImagesSection({ images, onRefresh }) {
  const [uploading, setUploading] = useState(false);

  const handleUpload = async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    const formData = new FormData();
    formData.append('image', file);
    formData.append('type', 'room');
    formData.append('room_type', prompt('Enter room type (standard, deluxe, suite, vip, premium, business):') || 'custom');

    setUploading(true);
    const token = localStorage.getItem('token');
    await fetch('/api/admin/images/upload', {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}` },
      body: formData
    });
    setUploading(false);
    onRefresh();
  };

  const handleDelete = async (filename) => {
    if (confirm('Delete image?')) {
      await api.delete('/api/admin/images/delete?filename=' + encodeURIComponent(filename));
      onRefresh();
    }
  };

  return (
    <div className="admin-section active">
      <h2>Image Management</h2>
      
      <div className="upload-zone" onClick={() => document.getElementById('imageInput').click()}>
        <input type="file" id="imageInput" accept="image/*" onChange={handleUpload} style={{ display: 'none' }} />
        {uploading ? '⏳ Uploading...' : '📁 Click to upload an image'}
      </div>

      <h3 style={{ marginTop: '1.5rem', marginBottom: '1rem' }}>Current Images</h3>
      <div className="image-grid">
        {images.map(img => (
          <div key={img.name} className="image-card">
            <img src={img.path} alt={img.name} />
            <button className="delete-btn" onClick={() => handleDelete(img.name)}>×</button>
            <div style={{ position: 'absolute', bottom: 0, left: 0, right: 0, background: 'rgba(0,0,0,0.7)', color: 'white', padding: '0.25rem', fontSize: '0.625rem' }}>
              {img.name}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

// Main Admin App
function AdminApp() {
  const [isAdmin, setIsAdmin] = useState(false);
  const [loading, setLoading] = useState(true);
  const [section, setSection] = useState('dashboard');
  const [stats, setStats] = useState({});
  const [rooms, setRooms] = useState(window.ADMIN_DATA?.rooms || []);
  const [roomTypes, setRoomTypes] = useState(window.ADMIN_DATA?.roomTypes || []);
  const [services, setServices] = useState(window.ADMIN_DATA?.services || []);
  const [packages, setPackages] = useState(window.ADMIN_DATA?.packages || []);
  const [bookings, setBookings] = useState([]);
  const [users, setUsers] = useState([]);
  const [reviews, setReviews] = useState([]);
  const [images, setImages] = useState([]);

  // Check access on load
  useEffect(() => {
    if (!checkAdmin()) {
      window.location.href = '/login';
      return;
    }
    setIsAdmin(true);
    setLoading(false);
  }, []);

  const loadData = async () => {
    if (!isAdmin) return;
    try {
      const [statsData, bookingsData, usersData, reviewsData, imagesData, roomsData, typesData, servicesData, packagesData] = await Promise.all([
        api.get('/api/admin/stats'),
        api.get('/api/admin/bookings'),
        api.get('/api/admin/users'),
        api.get('/api/reviews/list'),
        api.get('/api/admin/images'),
        api.get('/api/rooms/list'),
        api.get('/api/rooms/types/list'),
        api.get('/api/services/list'),
        api.get('/api/packages/list')
      ]);
      
      setStats(statsData);
      setBookings(bookingsData || []);
      setUsers(usersData || []);
      setReviews(reviewsData || []);
      setImages(imagesData || []);
      setRooms(roomsData || []);
      setRoomTypes(typesData || []);
      setServices(servicesData || []);
      setPackages(packagesData || []);
    } catch (err) {
      console.error('Error loading data:', err);
    }
  };

  useEffect(() => {
    if (isAdmin) {
      loadData();
    }
    
    // Navigation
    document.querySelectorAll('.admin-nav a[data-section]').forEach(link => {
      link.addEventListener('click', (e) => {
        e.preventDefault();
        document.querySelectorAll('.admin-nav a').forEach(l => l.classList.remove('active'));
        link.classList.add('active');
        setSection(link.dataset.section);
      });
    });
  }, [isAdmin]);

  if (loading) {
    return <div style={{padding: '2rem', textAlign: 'center'}}>Checking access...</div>;
  }

  if (!isAdmin) {
    return <div style={{padding: '2rem', textAlign: 'center', color: '#c00'}}>Access denied</div>;
  }

  return (
    <div>
      <div className="admin-header">
        <h1>{sectionTitles[section] || 'Admin Panel'}</h1>
        <button className="btn btn-outline btn-sm" onClick={loadData}>🔄 Refresh</button>
      </div>

      {section === 'dashboard' && <Dashboard stats={stats} />}
      {section === 'rooms' && <RoomsSection rooms={rooms} roomTypes={roomTypes} onRefresh={loadData} />}
      {section === 'room-types' && <RoomTypesSection roomTypes={roomTypes} onRefresh={loadData} />}
      {section === 'bookings' && <BookingsSection bookings={bookings} onRefresh={loadData} />}
      {section === 'calendar' && <CalendarSection bookings={bookings} />}
      {section === 'services' && <ServicesSection services={services} onRefresh={loadData} />}
      {section === 'packages' && <PackagesSection packages={packages} onRefresh={loadData} />}
      {section === 'users' && <UsersSection users={users} onRefresh={loadData} />}
      {section === 'reviews' && <ReviewsSection reviews={reviews} onRefresh={loadData} />}
      {section === 'images' && <ImagesSection images={images} onRefresh={loadData} />}
    </div>
  );
}

const sectionTitles = {
  'dashboard': '📊 Dashboard',
  'rooms': '🛏️ Rooms',
  'room-types': '📁 Room Types',
  'bookings': '📅 Bookings',
  'calendar': '🗓️ Calendar',
  'services': '🍽️ Services',
  'packages': '📦 Packages',
  'users': '👥 Users',
  'reviews': '⭐ Reviews',
  'images': '🖼️ Images'
};

ReactDOM.createRoot(document.getElementById('adminRoot')).render(<AdminApp />);
