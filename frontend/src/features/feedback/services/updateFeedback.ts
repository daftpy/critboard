import { FeedbackFormData } from "../types/feedbackTypes.ts";

export async function updateFeedback(formData: FeedbackFormData) {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_URL}/feedback/${formData.commentId}`,
      {
        method: "PATCH", // Using PATCH method for partial update
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(formData),
      },
    );

    const clonedResponse = response.clone();

    const rawResponse = await clonedResponse.text();
    console.log("Raw response: ", rawResponse);

    const data = await response.json();

    if (response.ok) {
      return {
        type: "success",
        message: data.message,
        feedback: data.feedback,
      };
    } else {
      return { type: "error", errors: data.errors };
    }
  } catch (error) {
    const errorMessage = (error as Error).message;
    console.log("There was a problem: ", errorMessage);
    return { type: "error", errors: [errorMessage] };
  }
}
