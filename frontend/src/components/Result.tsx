import { useState } from "react";

interface ResultProps {
    shortURL: string;
    onReset?: () => void;
}

export default function Result({ shortURL, onReset }: ResultProps) {

    const [copied, setCopied] = useState<boolean>(false);

    // copy to clipboard handler
    const handleCopy = async() => {
        try{
            await navigator.clipboard.writeText(shortURL);
            setCopied(true);

            // after 2 seconds, reset copied state
            setTimeout(() => setCopied(false), 2000);
        } catch (err) {
            console.error('Failed to copy:', err);
            alert('Failed to copy to clipboard');
        }
    }

    return (
        <div className="mt-8 pt-6 border-t-2 border-gray-200">
            <div className="text-center mb-6">
              <span className="inline-flex items-center px-4 py-2 bg-green-100 text-green-700 rounded-full font-semibold text-lg">
                 ✅ Short URL Generated Successfully!
              </span>
            </div>
            {/* Short URL Display */}
            <div className="flex items-center gap-3 p-4 bg-gray-50 rounded-xl border-2 border-gray-200 mb-4">
              <span className="flex-1 text-lg font-medium text-purple-600 break-all">
                {shortURL}
              </span>
              <button
                onClick={handleCopy}
                className={`px-5 py-2 rounded-lg font-medium transition-all whitespace-nowrap ${
                    copied
                        ? 'bg-green-500 text-white'
                        : 'bg-purple-600 text-white hover:bg-purple-700'
                }`}
              >
                {copied ? 'Copied!' : 'Copy'}
              </button>
            </div>
            {/*  Action Buttons */}
            <div className='grid grid-cols-1 md:grid-cols-2 gap-3'>
              <a
                href={shortURL}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center justify-center px-6 py-3 bg-green-500 text-white rounded-xl font-medium hover:bg-green-600 transition-all transform hover:scale-[1.02] shadow-md"
              >
              🔗 Visit Short URL
              </a>
              <button
                onClick={onReset}
                className="px-6 py-3 bg-orange-500 text-white rounded-xl font-medium hover:bg-orange-600 transition-all transform hover:scale-[1.02] shadow-md"
              >
                ↻ Generate New
              </button>
            </div>
        </div>
    );
}