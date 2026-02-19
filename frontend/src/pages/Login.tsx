import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import Layout from "../components/Layout";
import { useAuth } from "../auth/AuthContext";

export default function Login() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    setError("");
    setIsSubmitting(true);

    try {
      await login(email, password);
      navigate("/");
    } catch (err: any) {
      const message = err?.response?.data?.error || "Login failed";
      setError(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Layout>
      <section className="mx-auto max-w-md px-4 py-12">
        <div className="rounded-2xl border border-gray-200 bg-white p-6 shadow-sm">
          <h1 className="text-2xl font-semibold text-gray-900">Welcome back</h1>
          <p className="mt-1 text-sm text-gray-600">Log in to continue your analysis.</p>

          <form onSubmit={handleSubmit} className="mt-6 space-y-4">
            <label className="block">
              <span className="text-sm font-medium text-gray-700">Email</span>
              <input
                type="email"
                autoComplete="email"
                required
                value={email}
                onChange={(event) => setEmail(event.target.value)}
                className="mt-2 w-full rounded-lg border border-gray-300 p-3 text-sm text-gray-900 focus:border-[#55F] focus:outline-none"
                placeholder="you@example.com"
              />
            </label>

            <label className="block">
              <span className="text-sm font-medium text-gray-700">Password</span>
              <input
                type="password"
                autoComplete="current-password"
                required
                value={password}
                onChange={(event) => setPassword(event.target.value)}
                className="mt-2 w-full rounded-lg border border-gray-300 p-3 text-sm text-gray-900 focus:border-[#55F] focus:outline-none"
                placeholder="Your password"
              />
            </label>

            {error && (
              <div className="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">
                {error}
              </div>
            )}

            <button
              type="submit"
              disabled={isSubmitting}
              className="w-full rounded-lg bg-[#55F] px-4 py-2 text-sm font-semibold text-white hover:bg-[#44E] disabled:cursor-not-allowed disabled:opacity-70"
            >
              {isSubmitting ? "Signing in..." : "Login"}
            </button>
          </form>

          <p className="mt-4 text-sm text-gray-600">
            New here?{" "}
            <Link to="/register" className="font-semibold text-[#55F] hover:text-[#44E]">
              Create an account
            </Link>
          </p>
        </div>
      </section>
    </Layout>
  );
}
