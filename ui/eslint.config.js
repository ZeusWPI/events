import antfu from "@antfu/eslint-config";
import pluginQuery from "@tanstack/eslint-plugin-query";

export default antfu(
  {
    react: true,

    typescript: {
      tsconfigPath: "tsconfig.json",
    },

    stylistic: {
      indent: 2,
      quotes: "double",
      semi: true,
    },
    ignores: [
      "src/components/ui/*",
    ],
    rules: {
      "ts/strict-boolean-expressions": "off",
      "react-hooks-extra/no-direct-set-state-in-use-effect": "off",
    },
  },
).prepend(...pluginQuery.configs["flat/recommended"]);
