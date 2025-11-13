// import { useRef } from "react"
import { useState, useRef } from "react"
import { useMutation } from "@tanstack/react-query"
import { generateShortURL, type generateURLRequest, type generateURLResponse} from "../utils/http"
import type React from "react";


interface URLFormProps {
    setURLHandler: (url: string | null) => void;
}

export default function URLForm({ setURLHandler }: URLFormProps){

    const inputRef = useRef<HTMLInputElement>(null);
    const [error, setError] = useState<string>("");
    
    const { mutate: generateFn, isPending } = useMutation({
        mutationFn: ({ url } : generateURLRequest) => {
            return generateShortURL({ request: { url } })
        },
        onSuccess: ({ message, short_url }: generateURLResponse) => {
            console.log(message);
            setError("");
            setURLHandler(short_url);
        },
        onError: (e) => {
            console.error("Failed to generate short URL:", e);
            setError("Failed to generate short URL. Please try again.");
            setURLHandler(null);
        }
    })

    const submitHandler = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        const longURL = inputRef.current?.value || "";

        // clear previous error
        setError("");
        setURLHandler(null);

        // validate input
        if (!longURL){
            setError("Please enter a URL.");
            return;
        }

        // simple URL form validation
        try {
            new URL(longURL);
        } catch {
            setError("Please enter a valid URL.");
            return;
        }

        generateFn({ url: longURL });
    }

    return (
        <form className='space-y-4' onSubmit={submitHandler}>
          <div className='relative'>
            <input
              ref={inputRef}
              type="text"
              placeholder="Enter your long URL here (e.g., https://www.google.com)"
              className={`w-full px-5 py-4 text-stone-800 text-lg border-2 rounded-xl outline-none transition-all disabled:bg-gray-100 disabled:cursor-not-allowed
                  ${error
                      ? 'border-red-500 focus:border-red-500 focus:ring-4 focus:ring-red-200' 
                      : 'border-gray-300 focus:border-stone-500 focus:ring-4 focus:ring-stone-200'
                  }
              `}
            />
          </div>
          <button
            type="submit"
            disabled={isPending}
            className="w-full py-4 px-6 text-lg font-semibold text-white bg-linear-to-r from-purple-600 to-indigo-600 rounded-xl hover:from-purple-700 hover:to-indigo-700 focus:outline-none focus:ring-4 focus:ring-purple-300 disabled:opacity-60 disabled:cursor-not-allowed transform hover:scale-[1.02] transition-all duration-200 shadow-lg"
          >
            {isPending ? (
                <span className="flex items-center justify-center">
                    <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    Generating...
                </span>
            ) : (
                'Generate'
            )}
          </button>
        </form>
    )
}