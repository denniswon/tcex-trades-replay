import * as yup from "yup";

const schema = yup.object().shape({
  replay_rate: yup
    .string()
    .matches(/^\d+\.?\d*$/, "Must be a non-negative number")
    .required("Replay rate is required"),
});

export default schema;
