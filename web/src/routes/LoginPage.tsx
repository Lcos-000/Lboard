export default function LoginPage() {
  return (
    <section className="page-card">
      <h1>登录</h1>
      <p>Phase 0：这里只是登录页面占位，Phase 1 再实现真实登录。</p>

      <form className="form">
        <input placeholder="邮箱" disabled />
        <input placeholder="密码" type="password" disabled />
        <button type="button" disabled>
          登录
        </button>
      </form>
    </section>
  );
}
