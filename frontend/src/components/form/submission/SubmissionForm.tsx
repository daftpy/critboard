import { useNavigate } from "react-router-dom";
import styles from "../../../styles/components/form/submission/SubmissionForm.module.css";
import Button from "../../ui/Button";
import TextInput from "../../ui/TextInput";
import TypeSelector from "./TypeSelector";
import { IFormData, useSubmissionForm } from "./useSubmissionForm";

export default function SubmissionForm() {
  const navigate = useNavigate();

  const inititalData: IFormData = {
    title: "",
    description: "",
    type: "LINK",
    link: "",
  };

  const { formData, currentType, changeType, handleChange, handleSubmit } =
    useSubmissionForm(inititalData, navigate);

  const typeField =
    currentType === "LINK" ? (
      <>
        <label>Link</label>
        <TextInput
          name="link"
          value={formData.link}
          placeholder="https://example.com"
          onChange={handleChange}
        />
      </>
    ) : (
      <input type="file" name="uploadfile" />
    );

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
      <div className={styles.inputField}>
        <label>Title</label>
        <TextInput
          name="title"
          value={formData.title}
          placeholder="Enter a title"
          onChange={handleChange}
        />
      </div>
      <div className={styles.inputField}>
        <label>Description</label>
        <textarea
          onChange={handleChange}
          className={styles.textArea}
          name="description"
          value={formData.description}
          placeholder="Enter your description"
        />
      </div>
      <div>
        <TypeSelector selectType={changeType} currentType={currentType} />
        <div className={`${styles.typeField} ${styles.inputField}`}>
          {typeField}
        </div>
      </div>
      <Button message="Submit" />
    </form>
  );
}
