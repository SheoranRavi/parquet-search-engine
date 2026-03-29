export default function SearchMeta({ totalResults, queryTime }) {
  return (
    <div className="text-sm text-gray-500">
      {totalResults} result{totalResults !== 1 && "s"} in {queryTime}ms
    </div>
  );
}