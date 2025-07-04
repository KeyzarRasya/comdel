import type { PageServerLoad } from "../../dashboard/$types";

export const load: PageServerLoad = async ({url}) => {
    const statusCode: string = url.searchParams.get("status_code") || ""
    const transactionStatus: string = url.searchParams.get("transaction_status") || ""
    const orderId: string = url.searchParams.get("order_id") || ""

    return {
        statusCode,
        transactionStatus,
        orderId,
    }
}