<template>
  <div class="form-container card">
    <h2>Yeni İşlem Ekle</h2>

    <form @submit.prevent="submitTransaction">

      <label>Varlık</label>
      <select v-model="form.asset" required>
        <option disabled value="">Seçiniz</option>
        <option value="try">TRY</option>
        <option value="usd">USD</option>
        <option value="eur">EUR</option>
        <option value="gbp">GBP</option>
        <option value="chf">CHF</option>
        <option value="ceyrek-altin">Çeyrek Altın</option>
        <option value="yarim-altin">Yarım Altın</option>
        <option value="tam-altin">Tam Altın</option>
        <option value="cumhuriyet-altini">Cumhuriyet Altını</option>
        <option value="gram-altin">Gram Altın</option>
        <option value="gumus">Gümüş</option>
      </select>

      <label>İşlem Tipi</label>
      <select v-model="form.type" required>
        <option value="add">Ekle</option>
        <option value="subtract">Çıkar</option>
      </select>

      <label>Miktar</label>
      <input type="number" v-model="form.amount" step="0.01" required />

      <label>İşlem Tarihi</label>
      <input type="datetime-local" v-model="form.transaction_date" />

      <label>Açıklama</label>
      <input type="text" v-model="form.description" />

      <button type="submit">Kaydet</button>

    </form>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import api from "../api";

const router = useRouter();

const form = ref({
  type: "add",
  asset: "",
  amount: "",
  transaction_date: "",
  description: ""
});

const submitTransaction = async () => {
  try {

    let payload = { ...form.value };

    if (payload.transaction_date) {
      payload.transaction_date = new Date(payload.transaction_date).toISOString();
    }

    await api.post("/transactions", payload);

    alert("Transaction eklendi!");

    router.push("/dashboard");

  } catch (err) {
    console.error(err);
    alert("İşlem eklenemedi");
  }
};
</script>

<style scoped>
.form-container {
  max-width: 500px;
  margin: 2rem auto;
}
label {
  display: block;
  margin-bottom: 0.4rem;
  margin-top: 1rem;
}
</style>