import { redirect } from "@sveltejs/kit";
import type {PageServerLoad } from "./$types";
import { VITE_DEV_ENV, VITE_DOCKER_URL, VITE_LOCAL_URL } from "$env/static/private";

export const load: PageServerLoad = async ({cookies, fetch}) => {
    const token = cookies.get("jwt");

    const BASE_URL = VITE_DEV_ENV == "dev" ? VITE_LOCAL_URL : VITE_DOCKER_URL

    console.log(BASE_URL)

    
    if (token === undefined) {
        throw redirect(307, "/login")
    }
    
    try {
        const response = await fetch(`${BASE_URL}/user/info`, {
            headers: {
                Cookie:`jwt=${token}`
            },
            method:"GET"
        })
        
        console.log("i am")
        console.log(response.ok)
    
            const userInfo = await response.json();
            console.log(userInfo)
    
        return {
            token,
            userInfo,
            baseUrl:BASE_URL
        }

    }catch(err){
        console.log(err)
        throw redirect(307, "/error")
    }
    
}