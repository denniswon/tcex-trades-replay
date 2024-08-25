import * as yup from "yup";

const schema = yup.object().shape({
  replay_rate: yup
    .string()
    .matches(/^\d+\.?\d*$/, "Must be a non-negative number")
    .nullable()
    .default("0.0"),
  filename: yup.string().default("trades.txt").nullable(),
});

export default schema;
