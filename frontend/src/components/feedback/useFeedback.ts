import { useState } from "react";
import { getReplies } from "../../services/feedback/getFeedback";
import { FeedbackData } from "./Feedback";
import { useFeedbackData } from "./useFeedbackData";
import { ReplyButtonProps } from "../ui/feedback/ReplyButton";
import { MetaProps } from "./FeedbackMeta";
import { DeleteButtonProps } from "../ui/feedback/DeleteButton";
import { EditButtonProps } from "../ui/feedback/EditButton";

type EditFormProps = {
  commentId: string;
  text: string;
  replyForm: false;
  buttonText: string;
  actionType: "POST" | "UPDATE";
  onSubmit: (updatedFeedback: FeedbackData) => void;
}

export function useFeedback(feedback: FeedbackData, updateFeedback: (updatedFeedback: FeedbackData) => void) {
  const [showReplies, setShowReplies] = useState<boolean>(false);
  const [showForm, setShowForm] = useState<boolean>(false);
  const [editMode, setEditMode] = useState<boolean>(false);
  const [deleteConfirm, setDeleteConfirm] = useState<boolean>(false);

  const {
    feedbackData,
    setFeedbackData,
    addFeedbackData,
    updateFeedbackData
  } = useFeedbackData();

  const toggleForm = () => {
    setShowForm(!showForm);
  }

  const toggleEditMode = () => {
    setEditMode(!editMode);
  }

  const toggleReplies = () => {
    async function fetchReplies() {
      const result = await getReplies(feedback.commentId);

      if (result.type === "success") {
        console.log(result);
        setFeedbackData(result.feedback || []);
      } else {
        console.error(result.errors);
      }
    }

    if (showReplies) {
      setShowReplies(false);
    } else {
      fetchReplies().then(() => {
        setShowReplies(true);
        if (feedback.replies === 0) {
          setShowForm(true);
        }
      })
    }
  }

  const addFeedback = (newFeedback: FeedbackData) => {
    if (editMode) {
      setEditMode(false);
      console.log('closing edit form');
    } else {
      feedback.replies += 1;
      addFeedbackData(newFeedback);
      if (!showReplies) {
        toggleReplies();
      }
      setShowForm(false);
    }
  }

  const setConfirm = (confirm: boolean) => {
    setDeleteConfirm(confirm);
  }

  const deleteFeedback = () => {
    console.log('Removed:', feedback.commentId);
  }

  const getEditProps = (): EditButtonProps => {
    return {
      editMode: editMode,
      onClick: toggleEditMode
    }
  }

  const getDeleteProps = (): DeleteButtonProps => {
    return {
      onClick: deleteFeedback,
      confirm: deleteConfirm,
      setConfirm: setConfirm
    }
  }

  const getReplyProps = (): ReplyButtonProps => {
    return {
      toggleReplies: toggleReplies,
      toggleForm: toggleForm,
      replyCount: feedback.replies,
      showForm: showForm,
      showReplies: showReplies
    }
  }

  const getMetaProps = (): MetaProps => {
    return {
      edit: getEditProps(),
      remove: getDeleteProps(),
      reply: getReplyProps(),
      createdAt: feedback.createdAt,
      author: "author"
    }
  }

  const getEditFormProps = (): EditFormProps => {
    return {
      commentId: feedback.commentId,
      text: feedback.feedbackText,
      replyForm: false,
      buttonText: "Save Edit",
      actionType: "UPDATE",
      onSubmit: updateFeedback
    }
  }

  return {
    showReplies,
    showForm,
    editMode,
    feedbackData,
    updateFeedbackData,
    toggleForm,
    toggleReplies,
    addFeedback,
    getMetaProps,
    getEditFormProps
  }
}