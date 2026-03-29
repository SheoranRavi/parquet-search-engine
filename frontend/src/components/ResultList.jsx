export default function ResultList({ results }) {
  if (results.length === 0) {
    return <p className="text-gray-400 text-center mt-8">No results found.</p>;
  }

  return (
    <ul className="space-y-3">
      {results.map((msg, i) => (
        <li key={i} className="p-4 border border-gray-200 rounded-lg">
          <p className="text-gray-800">{msg.Message}</p>
        </li>
      ))}
    </ul>
  );
}