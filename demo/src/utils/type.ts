type GetParameterType<T extends (...args: any) => any> = T extends (
  param: infer P
) => any
  ? P
  : never;

type ParameterType = GetParameterType<(value: string) => void>; // ParameterType is inferred to be string
