import { useEffect, useState } from 'preact/hooks'
import { Search, Filter } from 'lucide-preact';

type BinRequest = {
  Method: string;
  URL: string;
  Time: string;
  Headers: Record<string, string[]>;
  Body: string;
};

export const App = () => {
  const [binId, setBinId] = useState('');
  const [requests, setRequests] = useState<BinRequest[]>([]);
  const [selectedRequest, setSelectedRequest] = useState<BinRequest | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterMethod, setFilterMethod] = useState('');

  useEffect(() => {
    fetch('/api/bin')
      .then(response => response.json())
      .then(data => setBinId(data.bin));

    const ws = new WebSocket(`${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/ws`);
    ws.onmessage = (event) => {
      const newRequest = JSON.parse(event.data);
      setRequests(prevRequests => [newRequest, ...prevRequests]);
    };

    return () => ws.close();
  }, []);

  useEffect(() => {
    if (requests.length > 0) {
      window.document.title = `Requests - (${requests.length}) request${'s'.repeat(requests.length > 1 ? 1 : 0)}`;
    } else {
      window.document.title = 'Requests';
    }
  }, [requests]);

  const filteredRequests = requests.filter(request =>
    request.Method.toLowerCase().includes(filterMethod.toLowerCase()) &&
    (request.URL.toLowerCase().includes(searchTerm.toLowerCase()) ||
      request.Method.toLowerCase().includes(searchTerm.toLowerCase()))
  );

  return (
    <div className="flex h-screen bg-gray-100">
      <div className="w-1/3 bg-white p-6 overflow-auto">
        <h1 className="text-2xl font-bold mb-4">HTTP Request Bin</h1>
        {
          binId && (
            <>
              <p className="mb-4">Bin ID: {binId}</p>
              <p className="mb-4">Endpoint: {' '}
                <a href={`http://${binId}.r.ahh.bet`}>{binId}.r.ahh.bet</a>
              </p>
            </>
          )
        }

        <div className="mb-4 flex">
          <div className="relative flex-grow">
            <input
              type="text"
              placeholder="Search requests..."
              className="w-full p-2 pl-8 border rounded"
              value={searchTerm}
              onChange={(e) => setSearchTerm((e.target as HTMLInputElement).value)}
            />
            <Search className="absolute left-2 top-2 text-gray-400" size={20} />
          </div>
          <div className="relative ml-2">
            <select
              className="p-2 pl-8 border rounded appearance-none"
              value={filterMethod}
              onChange={(e) => setFilterMethod((e.target as HTMLSelectElement).value)}
            >
              <option value="">All Methods</option>
              <option value="GET">GET</option>
              <option value="POST">POST</option>
              <option value="PUT">PUT</option>
              <option value="DELETE">DELETE</option>
            </select>
            <Filter className="absolute left-2 top-2 text-gray-400" size={20} />
          </div>
        </div>

        {filteredRequests.map((request, index) => (
          <div
            key={index}
            className={`p-2 mb-2 rounded cursor-pointer ${selectedRequest === request ? 'bg-blue-100' : 'hover:bg-gray-100'
              }`}
            onClick={() => setSelectedRequest(request)}
          >
            <div className="font-bold">{request.Method} {request.URL}</div>
            <div className="text-sm text-gray-500">{request.Time}</div>
          </div>
        ))}
      </div>
      <div className="flex-1 p-6 bg-white m-6 rounded-lg shadow overflow-auto">
        {selectedRequest ? (
          <div>
            <h2 className="text-xl font-bold mb-4">Request Details</h2>
            <pre className="bg-gray-100 p-4 rounded overflow-auto">
              {JSON.stringify(selectedRequest, null, 2)}
            </pre>
          </div>
        ) : (
          <div className="text-center text-gray-500 mt-10">
            Select a request to view details
          </div>
        )}
      </div>
    </div>
  );
};

