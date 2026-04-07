import { Link } from "react-router-dom";
import Layout from "../components/Layout";

export default function Home() {
  return (
    <Layout>
      <div className="flex flex-1 flex-col sm:flex-row min-h-[calc(100vh-7rem)]">
        <Link
          to="/analysis"
          className="group flex flex-1 flex-col items-center justify-center gap-4 bg-indigo-600 px-8 py-16 text-center transition-colors hover:bg-indigo-500"
        >
          <h2 className="text-3xl font-bold text-white">Analysis</h2>
          <span className="mt-2 rounded-md border border-white px-5 py-2 text-sm font-semibold text-white group-hover:bg-white group-hover:text-indigo-600 transition-colors">
            Get started →
          </span>
        </Link>
        <Link
          to="/vocabulary"
          className="group flex flex-1 flex-col items-center justify-center gap-4 bg-gray-900 px-8 py-16 text-center transition-colors hover:bg-gray-800"
        >
          <h2 className="text-3xl font-bold text-white">Vocabulary</h2>
          <span className="mt-2 rounded-md border border-white px-5 py-2 text-sm font-semibold text-white group-hover:bg-white group-hover:text-gray-900 transition-colors">
            Browse words →
          </span>
        </Link>
      </div>
    </Layout>
  );
}          
