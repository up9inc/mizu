import { Backdrop, Box, Fade, Modal } from "@material-ui/core";
import { RedocStandalone } from "redoc";
import Api from "../../helpers/api";

const api = new Api();

const OasDModal = ({ openModal, handleCloseModal }) => {
  const getOASAllSpecs = async () => {
    const data = api.getOASAllSpecs();
  };

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
          <RedocStandalone 
            specUrl="http://petstore.swagger.io/v2/swagger.json"
            options={{
                theme: {colors: {primary:{main:"#fafafa",light:"#fafafa", dark:"#fafafa"}}}
            }} />
        </Box>
      </Fade>
    </Modal>
  );
};

export default OasDModal;
