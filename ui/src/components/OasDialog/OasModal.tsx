import { Backdrop, Box, Fade, Modal } from "@material-ui/core";
import { useEffect, useState } from "react";
import { RedocStandalone } from "redoc";
import Api from "../../helpers/api";

const api = new Api();

const OasDModal = ({ openModal, handleCloseModal, selectedService }) => {
  const [serviceOAS, setServiceOAS] = useState(null);

  useEffect(() => {
    (async () => {
      try {
        const data = await api.getOASAByService(selectedService);
        setServiceOAS(data);
      } catch (e) {
        console.error(e);
      }
    })();
  }, [selectedService]);

  return (
    <Modal
      aria-labelledby="transition-modal-title"
      aria-describedby="transition-modal-description"
      open={openModal}
      onClose={handleCloseModal}
      closeAfterTransition
      BackdropComponent={Backdrop}
      BackdropProps={{
        timeout: 500,
      }}
      style={{ overflow: "auto", backgroundColor: "#fafafa" }}
    >
      <Fade in={openModal}>
        <Box>
          {serviceOAS && (
            <RedocStandalone
              spec={serviceOAS}
              options={{
                theme: {
                  colors: {
                    primary: {
                      main: "#fafafa",
                      light: "#fafafa",
                      dark: "#fafafa",
                    },
                  },
                },
              }}
            />
          )}
        </Box>
      </Fade>
    </Modal>
  );
};

export default OasDModal;
