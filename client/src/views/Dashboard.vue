<template>
  <div class="dashboard-container">
    <!-- Portfolio Total -->
    <div class="portfolio-total card">
      <h2>Toplam Varlık</h2>
      <p>{{ totalPrice.toFixed(2) }} {{ currency.toUpperCase() }}</p>
    </div>

    <!-- Chart -->
    <div class="portfolio-chart card">
      <h2>Mal Varlığı Grafiği</h2>
      <select v-model="chartType">
        <option value="pie">Pasta</option>
        <option value="bar">Çubuk</option>
      </select>
      <canvas ref="chartCanvas"></canvas>
    </div>

    <!-- Asset List -->
    <div class="asset-section">
      <h2>Döviz</h2>
      <div class="asset-list">
        <div v-for="asset in groupedAssets.currency" :key="asset.id" class="asset-card" @click="toggleDetails(asset.id)">
          <div class="asset-header">
            <span class="asset-name">{{ asset.asset.toUpperCase() }}</span>
            <span class="asset-total">{{ asset.total_price_by_asset.toFixed(2) }} {{ currency.toUpperCase() }}</span>
          </div>
          <div v-if="expandedAssetId === asset.id" class="asset-details">
            <p><strong>Miktar:</strong> {{ asset.amount }}</p>
            <p><strong>Fiyat:</strong> {{ asset.price.toFixed(2) }} {{ currency.toUpperCase() }}</p>
            <p><strong>Son İşlemler:</strong></p>
            <ul>
              <li v-for="tx in asset.transactions" :key="tx.id">
                {{ tx.type.toUpperCase() }} {{ tx.amount }} @ {{ tx.price.toFixed(2) }} {{ currency.toUpperCase() }} → Toplam: {{ tx.total_price.toFixed(2) }} {{ currency.toUpperCase() }}
              </li>
            </ul>
          </div>
        </div>
      </div>

      <h2>Altın</h2>
      <div class="asset-list">
        <div v-for="asset in groupedAssets.gold" :key="asset.id" class="asset-card" @click="toggleDetails(asset.id)">
          <div class="asset-header">
            <span class="asset-name">{{ asset.asset.toUpperCase() }}</span>
            <span class="asset-total">{{ asset.total_price_by_asset.toFixed(2) }} {{ currency.toUpperCase() }}</span>
          </div>
          <div v-if="expandedAssetId === asset.id" class="asset-details">
            <p><strong>Miktar:</strong> {{ asset.amount }}</p>
            <p><strong>Fiyat:</strong> {{ asset.price.toFixed(2) }} {{ currency.toUpperCase() }}</p>
          </div>
        </div>
      </div>

      <h2>Gümüş</h2>
      <div class="asset-list">
        <div v-for="asset in groupedAssets.silver" :key="asset.id" class="asset-card" @click="toggleDetails(asset.id)">
          <div class="asset-header">
            <span class="asset-name">{{ asset.asset.toUpperCase() }}</span>
            <span class="asset-total">{{ asset.total_price_by_asset.toFixed(2) }} {{ currency.toUpperCase() }}</span>
          </div>
          <div v-if="expandedAssetId === asset.id" class="asset-details">
            <p><strong>Miktar:</strong> {{ asset.amount }}</p>
            <p><strong>Fiyat:</strong> {{ asset.price.toFixed(2) }} {{ currency.toUpperCase() }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from "vue";
import { Chart, registerables } from "chart.js";
import api from "../api";

Chart.register(...registerables);

const chartCanvas = ref(null);
const chartInstance = ref(null);
const chartType = ref("pie");

const assets = ref([]);
const groupedAssets = ref({ currency: [], gold: [], silver: [] });
const expandedAssetId = ref(null);

const totalPrice = ref(0);
const currency = ref("try"); // default

const chartData = ref({
  labels: [],
  datasets: [
    {
      label: "Mal Varlığı",
      data: [],
      backgroundColor: ["#42b883", "#ff6384", "#36a2eb", "#ffce56", "#9b59b6"],
    },
  ],
});

const fetchAssets = async () => {
  try {
    const res = await api.get("/assets/all");
    assets.value = res.data.data.assets || [];
    totalPrice.value = res.data.data.total_price || 0;
    currency.value = res.data.data.currency || "try";

    groupedAssets.value = categorizeAssets(assets.value);

    chartData.value.labels = assets.value.map(a => a.asset.toUpperCase());
    chartData.value.datasets[0].data = assets.value.map(a => a.total_price_by_asset);
    renderChart();
  } catch (err) {
    console.error("Assets fetch error:", err);
  }
};

const renderChart = () => {
  if (chartInstance.value) chartInstance.value.destroy();

  chartInstance.value = new Chart(chartCanvas.value, {
    type: chartType.value,
    data: chartData.value,
    options: {
      responsive: true,
      plugins: {
        legend: { position: "bottom" },
      },
    },
  });
};

const toggleDetails = (id) => {
  expandedAssetId.value = expandedAssetId.value === id ? null : id;
};

const categorizeAssets = (assets) => {
  const groups = { currency: [], gold: [], silver: [] };
  assets.forEach(a => {
    if (["try","usd","eur","gbp","chf"].includes(a.asset)) groups.currency.push(a);
    else if (a.asset.includes("altin") || a.asset.includes("bilezik")) groups.gold.push(a);
    else if (a.asset === "gumus") groups.silver.push(a);
  });
  return groups;
};

onMounted(fetchAssets);
watch(chartType, () => renderChart());
</script>

<style scoped>
.dashboard-container { display: flex; flex-direction: column; gap: 2rem; }

.portfolio-total, .portfolio-chart { max-width: 700px; }

.asset-section h2 { margin-top: 2rem; margin-bottom: 1rem; }

.asset-list { display: flex; flex-wrap: wrap; gap: 1rem; }

.asset-card { background-color: var(--color-surface); border: 1px solid var(--color-surface-light); border-radius: 12px; padding: 1rem; cursor: pointer; width: 250px; transition: 0.2s; }

.asset-card:hover { box-shadow: 0 4px 12px rgba(0,0,0,0.2); }

.asset-header { display: flex; justify-content: space-between; font-weight: 600; margin-bottom: 0.5rem; }

.asset-details { font-size: 0.9rem; margin-top: 0.5rem; color: var(--color-text); }
</style>