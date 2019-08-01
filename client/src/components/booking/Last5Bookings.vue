<template>
  <div>
    <div class="list is-hoverable">
      <div v-for="entry in bookings" :key="entry" class="list-item has-background-white">
        <div class="level">
          <div class="level-left">
            <h6 class="subtitle is-6">{{printDateTime(entry.TimeStamp)}}</h6>
            <br />
            <h6
              class="subtitle is-6"
              v-if="entry.ItemID === 0"
            >{{displayUserName(getUserByID(entry.UserID))}} ~ Payment: {{entry.TotalPrice * -1}}€</h6>
            <h6
              class="subtitle is-6"
              v-else
            >{{displayUserName(getUserByID(entry.UserID))}} ~ {{entry.Amount}}x {{displayItem(getItemByID(entry.ItemID))}}</h6>
          </div>
          <!-- TODO timeout -->
          <button class="button is-small" @click="undo(entry)">
            <span class="icon is-small">
              <font-awesome-icon icon="redo" />
            </span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  computed: {
    bookings() {
      return this.$store.state.last5Bookings;
    }
  },
  methods: {
    undo(bookEntry) {
      this.$http
        .post("deleteBookEntry", bookEntry)
        .then(() => {
          this.$store.commit("getLast5Bookings");
          var message = "".concat(
            "Deleted book entry from ",
            this.displayUserName(this.getUserByID(bookEntry.UserID))
          );
          this.$responseEventBus.$emit("successMessage", message);
        })
        .catch(() => {
          this.$responseEventBus.$emit(
            "failureMessage",
            "Couldn't undo book entry"
          );
        });
    }
  }
};
</script>