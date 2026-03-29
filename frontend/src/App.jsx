import { useState } from "react";
import SearchBar from "./components/SearchBar";
import SearchMeta from "./components/SearchMeta";
import ResultList from "./components/ResultList";

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:9080/api';

export default function App() {
  const [results, setResults] = useState([]);
  const [meta, setMeta] = useState(null);

  const handleSearch = async (query) => {
    const reqBody = {
      query: query
    };
    
    const res = await fetch(`${API_URL}/search`, {
      method: 'POST',
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(reqBody)
    });
    const data = await res.json();

    setResults(data.messages || []);
    setMeta({ totalResults: data.totalCount || 0, queryTime: data.duration });
  };

  return (
    <div className="max-w-2xl mx-auto p-6 space-y-4">
      <SearchBar onSearch={handleSearch} />
      {meta && <SearchMeta totalResults={meta.totalResults} queryTime={meta.queryTime} />}
      {meta && <ResultList results={results} />}
    </div>
  );
}