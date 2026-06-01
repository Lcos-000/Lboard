export default function RegisterPage() {
  return (
    <section className="page-card">
      <h1>注册</h1>
      <p>Phase 0：这里只是注册页面占位，Phase 1 再实现真实注册。</p>

      <form className="form">
        <input placeholder="用户名" disabled />
        <input placeholder="邮箱" disabled />
        <input placeholder="密码" type="password" disabled />
        <button type="button" disabled>
          注册
        </button>
      </form>
    </section>
  );
}
