import {computed, watch} from 'vue';
import {useCompanyStore} from "@/stores/modules/company.js";

export function usePublicNetworkCheck(init = null) {
  const companyStore = useCompanyStore()
  const isPublicNetwork = computed(() => {
    return companyStore?.is_public_network == 1
  })

  if (typeof init === "function") {
    watch(() => isPublicNetwork.value, (val) => {
      val && init()
    }, {
      immediate: true
    })
  }

  return {
    isPublicNetwork,
  }
}
