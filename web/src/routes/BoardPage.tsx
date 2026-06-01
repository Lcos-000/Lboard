import { useParams } from 'react-router-dom';

export default function BoardPage() {
  const { roomId } = useParams();

  return (
    <section className="board-page">
      <aside className="toolbar">
        <h2>工具栏</h2>
        <button disabled>选择</button>
        <button disabled>矩形</button>
        <button disabled>画笔</button>
        <button disabled>文本</button>
      </aside>

      <div className="board-placeholder">
        <h1>白板页面占位</h1>
        <p>Room ID: {roomId}</p>
        <p>Phase 0：这里只验证页面和路由，Canvas 引擎在后续 Phase 实现。</p>
      </div>
    </section>
  );
}
