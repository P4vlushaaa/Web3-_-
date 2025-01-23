import React, { useState } from 'react';

function App() {
  const [files, setFiles] = useState<FileList | null>(null);
  const [uploadResult, setUploadResult] = useState<any>(null);

  async function handleUpload() {
    if (!files || files.length === 0) return;
    const formData = new FormData();
    formData.append('file', files[0]);
    const resp = await fetch('http://localhost:8081/upload', {
      method: 'POST',
      body: formData
    });
    const json = await resp.json();
    setUploadResult(json);
  }

  return (
    <div>
      <h1>Web3 OnlyFans</h1>
      <input type="file" onChange={(e) => setFiles(e.target.files)} />
      <button onClick={handleUpload}>Upload to FrostFS</button>

      {uploadResult && (
        <div>
          <p>ObjectID: {uploadResult.object_id}</p>
        </div>
      )}
    </div>
  );
}

export default App;
