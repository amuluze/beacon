/**
 * @Author     : Amu
 * @Date       : 2024/12/10 16:47
 * @Description:
 */

declare global {
    interface Navigator {
        msSaveOrOpenBlob: (blob: Blob, fileName: string) => void
        browserLanguage: string
    }
}

export {}
