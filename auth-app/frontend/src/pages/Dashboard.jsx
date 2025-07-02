import React, { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import './Dashboard.css';

const Dashboard = () => {
  const { user, logout, updateProfile } = useAuth();
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState({
    username: user?.username || '',
    first_name: user?.first_name || '',
    last_name: user?.last_name || '',
  });
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setMessage('');
    setError('');

    const result = await updateProfile(formData);
    
    if (result.success) {
      setMessage('Profile updated successfully!');
      setIsEditing(false);
    } else {
      setError(result.error);
    }
    
    setLoading(false);
  };

  const handleCancel = () => {
    setFormData({
      username: user?.username || '',
      first_name: user?.first_name || '',
      last_name: user?.last_name || '',
    });
    setIsEditing(false);
    setError('');
    setMessage('');
  };

  const handleLogout = () => {
    logout();
  };

  return (
    <div className="dashboard-container">
      <div className="dashboard-header">
        <h1>Welcome to Your Dashboard</h1>
        <button onClick={handleLogout} className="logout-button">
          Logout
        </button>
      </div>

      <div className="dashboard-content">
        <div className="profile-card">
          <h2>Profile Information</h2>
          
          {message && <div className="success-message">{message}</div>}
          {error && <div className="error-message">{error}</div>}

          {!isEditing ? (
            <div className="profile-info">
              <div className="info-group">
                <label>Username:</label>
                <p>{user?.username}</p>
              </div>
              <div className="info-group">
                <label>Email:</label>
                <p>{user?.email}</p>
              </div>
              <div className="info-group">
                <label>First Name:</label>
                <p>{user?.first_name || 'Not provided'}</p>
              </div>
              <div className="info-group">
                <label>Last Name:</label>
                <p>{user?.last_name || 'Not provided'}</p>
              </div>
              <div className="info-group">
                <label>Account ID:</label>
                <p>{user?.id}</p>
              </div>
              
              <button
                onClick={() => setIsEditing(true)}
                className="edit-button"
              >
                Edit Profile
              </button>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="edit-form">
              <div className="form-group">
                <label htmlFor="username">Username</label>
                <input
                  type="text"
                  id="username"
                  name="username"
                  value={formData.username}
                  onChange={handleChange}
                  required
                  minLength="3"
                  maxLength="30"
                  className="form-input"
                />
              </div>

              <div className="form-group">
                <label htmlFor="first_name">First Name</label>
                <input
                  type="text"
                  id="first_name"
                  name="first_name"
                  value={formData.first_name}
                  onChange={handleChange}
                  className="form-input"
                />
              </div>

              <div className="form-group">
                <label htmlFor="last_name">Last Name</label>
                <input
                  type="text"
                  id="last_name"
                  name="last_name"
                  value={formData.last_name}
                  onChange={handleChange}
                  className="form-input"
                />
              </div>

              <div className="form-actions">
                <button
                  type="submit"
                  disabled={loading}
                  className="save-button"
                >
                  {loading ? 'Saving...' : 'Save Changes'}
                </button>
                <button
                  type="button"
                  onClick={handleCancel}
                  className="cancel-button"
                >
                  Cancel
                </button>
              </div>
            </form>
          )}
        </div>

        <div className="info-card">
          <h3>Quick Stats</h3>
          <div className="stats-grid">
            <div className="stat-item">
              <span className="stat-label">Member Since</span>
              <span className="stat-value">
                {user?.created_at ? new Date(user.created_at).toLocaleDateString() : 'N/A'}
              </span>
            </div>
            <div className="stat-item">
              <span className="stat-label">Last Updated</span>
              <span className="stat-value">
                {user?.updated_at ? new Date(user.updated_at).toLocaleDateString() : 'N/A'}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;