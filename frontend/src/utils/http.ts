import axios from "axios";
import { QueryClient } from "@tanstack/react-query";

export const queryClient = new QueryClient();

export type generateURLRequest = {
    url: string
}

export type generateURLResponse = {
    message: string,
    short_url: string,
}

export async function generateShortURL({ request }: { request: generateURLRequest }){

    const response = await axios.post(
        "http://localhost:80/v1/shorten",
        {
            url: request.url
        },
        {
            headers: {
                "Content-Type": "application/json"
            }
        },
    )

    if (response.status !== 201){
        const error = new Error("An error occurred while generating short URL");
        throw error
    }

    return response.data;
}