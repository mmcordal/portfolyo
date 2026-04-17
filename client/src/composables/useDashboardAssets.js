import { ref } from 'vue'
import { assetService, downloadBlob } from '../services/portfolio'
import { createStatusState, setError, setOk } from './dashboardShared'

export function useDashboardAssets(currency) {
    const assetsAll = ref(null)
    const singleAsset = ref(null)
    const status = createStatusState()

    async function fetchAssets() {
        try {
            const res = await assetService.getAll(currency.value)
            assetsAll.value = res.data
        } catch (err) {
            setError(status, err)
        }
    }

    async function fetchSingleAsset(asset) {
        try {
            const res = await assetService.getByAsset(asset, currency.value)
            singleAsset.value = res.data
            setOk(status, `${asset.toUpperCase()} detayı getirildi.`)
        } catch (err) {
            setError(status, err)
        }
    }

    async function downloadAssetsPdf() {
        try {
            const blob = await assetService.downloadPDF(currency.value)
            downloadBlob(blob, `portfolio-${currency.value}.pdf`)
            setOk(status, 'Portföy PDF indirildi.')
        } catch (err) {
            setError(status, err)
        }
    }

    return {
        currency,
        assetsAll,
        singleAsset,
        status,
        fetchAssets,
        fetchSingleAsset,
        downloadAssetsPdf,
    }
}