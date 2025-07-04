
import { redirect } from "@sveltejs/kit";
import type {PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({params, fetch, cookies}) => {
    const token = cookies.get("jwt");
    console.log(params.slug)
    
    try{
        const response = await fetch(`http://backend:8080/videos/information/?id=${params.slug}`, {
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

