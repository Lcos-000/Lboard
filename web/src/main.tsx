import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';

import App from './App';
import './styles.css';

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    {/* 这里支持路由导航 */}
    <BrowserRouter>
      {/* 这里是应用组件，负责渲染路由组件 */}
      <App />
    </BrowserRouter>
  </React.StrictMode>
);
