import {ConfigEnv, defineConfig, UserConfig} from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vitejs.dev/config/
export default defineConfig((env: ConfigEnv) => {
    let config: UserConfig = {
        plugins: [svelte()],
    }
    if (env.command === 'serve') {
        config.server = {
            proxy: {
                '/api': 'http://127.0.0.1:8080',
            }
        }
    }

    return config
})