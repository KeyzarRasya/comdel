
import { redirect } from "@sveltejs/kit";
import type {PageServerLoad } from "./$types";
import { VITE_DOCKER_URL, VITE_DEV_ENV, VITE_LOCAL_URL } from "$env/static/private";

export const load: PageServerLoad = async ({params, fetch, cookies}) => {
    const token = cookies.get("jwt");
    const BASE_URL = VITE_DEV_ENV == "dev" ? VITE_LOCAL_URL : VITE_DOCKER_URL
    
    try{
        const response = await fetch(`${BASE_URL}/videos/information/?id=${params.slug}`, {
            method:"GET",
            headers: {
                Cookie:`jwt=${token}`
            },
        });
    
        const result = await response.json();
    
        if (result.status === 403) {
            throw redirect(403, "http://backend:5173/dashboard");
        }
    
        console.log(params.slug);
        console.log(result)    
    
        return {result}
    }catch(err) {
        console.error(err);
        throw redirect(307, "/error")
    }
}

