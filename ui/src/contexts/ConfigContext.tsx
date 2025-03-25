import { createContext, useContext } from 'react';

export type Config = {
  apiURL: string;
  basePath: string;
  title: string;
  navbarColor: string;
  version: string;  
  remoteNodes: string;
};

export const ConfigContext = createContext<Config>(null!);

export function useConfig() {
  return useContext(ConfigContext);
}
