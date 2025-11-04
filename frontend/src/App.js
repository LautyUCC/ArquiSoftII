import React, { useEffect, useState } from 'react';
import './App.css';

function App() {
  const [apiStatus, setApiStatus] = useState({
    users: 'checking...',
    properties: 'checking...',
    search: 'checking...'
  });

  useEffect(() => {
    // Verificar estado de las APIs
    const checkAPIs = async () => {
      try {
        const usersRes = await fetch('http://localhost:8080/health');
        const usersData = await usersRes.json();
        setApiStatus(prev => ({ ...prev, users: usersData.status }));
      } catch (err) {
        setApiStatus(prev => ({ ...prev, users: 'offline' }));
      }

      try {
        const propsRes = await fetch('http://localhost:8081/health');
        const propsData = await propsRes.json();
        setApiStatus(prev => ({ ...prev, properties: propsData.status }));
      } catch (err) {
        setApiStatus(prev => ({ ...prev, properties: 'offline' }));
      }

      try {
        const searchRes = await fetch('http://localhost:8082/health');
        const searchData = await searchRes.json();
        setApiStatus(prev => ({ ...prev, search: searchData.status }));
      } catch (err) {
        setApiStatus(prev => ({ ...prev, search: 'offline' }));
      }
    };

    checkAPIs();
  }, []);

  return (
    <div className="App">
      <header className="App-header">
        <h1>🏠 Spotly</h1>
        <p>Welcome to Spotly - Your perfect stay awaits</p>
        
        <div className="api-status">
          <h2>API Status</h2>
          <div className="status-grid">
            <div className={`status-card ${apiStatus.users === 'healthy' ? 'online' : 'offline'}`}>
              <h3>Users API</h3>
              <p>{apiStatus.users}</p>
            </div>
            <div className={`status-card ${apiStatus.properties === 'healthy' ? 'online' : 'offline'}`}>
              <h3>Properties API</h3>
              <p>{apiStatus.properties}</p>
            </div>
            <div className={`status-card ${apiStatus.search === 'healthy' ? 'online' : 'offline'}`}>
              <h3>Search API</h3>
              <p>{apiStatus.search}</p>
            </div>
          </div>
        </div>
      </header>
    </div>
  );
}

export default App;
