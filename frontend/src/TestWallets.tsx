import React, { useState, useEffect } from 'react'

const TestWallets: React.FC = () => {
  const [data, setData] = useState<any>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetch('/api/v1/wallets?userID=demo-user-123')
      .then(response => response.json())
      .then(data => {
        console.log('Raw API response:', data)
        setData(data)
      })
      .catch(err => {
        console.error('API error:', err)
        setError(err.message)
      })
  }, [])

  return (
    <div style={{ padding: '20px' }}>
      <h1>錢包資料測試</h1>
      {error && <div style={{ color: 'red' }}>錯誤: {error}</div>}
      {data ? (
        <div>
          <h2>原始 API 響應:</h2>
          <pre style={{ background: '#f5f5f5', padding: '10px', overflow: 'auto' }}>
            {JSON.stringify(data, null, 2)}
          </pre>
          
          <h2>錢包列表:</h2>
          {data.data?.data ? (
            <ul>
              {data.data.data.map((wallet: any) => (
                <li key={wallet.id}>
                  {wallet.name} - {wallet.type} - 餘額: {wallet.balance.amount} {wallet.balance.currency}
                </li>
              ))}
            </ul>
          ) : (
            <div>無法找到錢包資料</div>
          )}
        </div>
      ) : (
        <div>載入中...</div>
      )}
    </div>
  )
}

export default TestWallets