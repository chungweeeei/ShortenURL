import './App.css'

import { useState } from 'react'

import Header from './components/Header'
import URLForm from './components/Form'
import Footer from './components/Footer'
import Result from './components/Result'

function App() {

  const [shortURL, setShortURL] = useState<string | null>(null);
  const setURLHandler = (url: string | null) => {
    setShortURL(url);
  }

  const resetHandler = () => {
    setShortURL(null);
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <div className='w-full max-w-2xl'>

        {/* Main Card */}
        <div className='bg-white rounded-3xl shadow-2xl p-8 md:p-12'>

          {/* Header */}
          <Header />

          {/* URL Form */}
          <URLForm setURLHandler={setURLHandler}/>

          {/* Success Result */}
          {shortURL && <Result shortURL={shortURL} onReset={resetHandler}/>}
        </div>

        {/* Footer */}
        <Footer />
  
      </div>
    </div>
  )
}

export default App
